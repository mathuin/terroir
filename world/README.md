# World notes

## Level files

The only file really required for a directory to be a Minecraft level!

### Example

```
TAG_Compound(u''): 1 entries
{
  TAG_Compound(u'Data'): 23 entries
  {
    TAG_Byte(u'raining'): 0
    TAG_Int(u'generatorVersion'): 1
    TAG_Long(u'Time'): 4950
    TAG_Int(u'GameType'): 0
    TAG_String(u'generatorOptions'): u''
    TAG_Byte(u'MapFeatures'): 1
    TAG_String(u'generatorName'): u'default'
    TAG_Byte(u'initialized'): 1
    TAG_Byte(u'hardcore'): 0
    TAG_Long(u'RandomSeed'): 2603059821051629081
    TAG_Long(u'SizeOnDisk'): 0
    TAG_Byte(u'allowCommands'): 0
    TAG_Int(u'SpawnZ'): 220
    TAG_Long(u'LastPlayed'): 1406080954413
    TAG_Long(u'DayTime'): 4950
    TAG_Compound(u'GameRules'): 9 entries
    {
      TAG_String(u'doTileDrops'): u'true'
      TAG_String(u'doMobSpawning'): u'true'
      TAG_String(u'keepInventory'): u'false'
      TAG_String(u'naturalRegeneration'): u'true'
      TAG_String(u'doDaylightCycle'): u'true'
      TAG_String(u'doMobLoot'): u'true'
      TAG_String(u'mobGriefing'): u'true'
      TAG_String(u'doFireTick'): u'true'
      TAG_String(u'commandBlockOutput'): u'true'
    }
    TAG_Int(u'SpawnY'): 64
    TAG_Int(u'SpawnX'): -208
    TAG_Int(u'thunderTime'): 98924
    TAG_Int(u'version'): 19133
    TAG_Int(u'rainTime'): 17092
    TAG_Byte(u'thundering'): 0
    TAG_String(u'LevelName'): u'ExampleWorld'
  }
}
```

## Region files

### Coordinates

Region files have names like "region/r.x.z.mcr" where x and z are
"region coordinates".  These coordinates can be found this way:

```
regionX := chunkX >> 5
regionZ := chunkZ >> 5
```

### Format

Region files begin with 8kiB header with what chunks are present,
when last updated, and where found.

Location in file of chunk at chunk coordinates (x, z) can be found at:

```
xmod := Mod(x, 32)
if xmod < 0 {
  xmod = xmod + 32
}
zmod := Mod(z, 32)
if zmod < 0 {
  zmod = zmod + 32
}
offset := 4 * (xmod + zmod * 32)
```

Timestamp of chunk is 4096 bytes later in file.


| bytes     | description                |
| --------- | -------------------------- |
| 0-4095    | locations (1024 entries)   |
| 4096-8191 | timestamps (1024 entries)  |
| 8192+     | chunks and unused space?   |

## Chunks

### Location

| byte | description                          |
| ---- | ------------------------------------ |
|  0-2 | offset (big-endian, in 4KiB sectors) |
|   3  | count (number of 4KiB sectors)       |

Note: a chunk with an offset of 2 will begin right after timestamps table!

### Timestamps

| byte | description                              |
| ---- | ---------------------------------------- |
| 0-3  | timestamp (four-byte big-endian integer) |


### Data

| byte | description                              |
| ---- | ---------------------------------------- |
| 0-3  | length (in bytes) |
|   4  | compression type (1=gzip, 2=zlib) |
|   5  | compressed data (length-1 bytes) |

Note: all chunks must be padded to multiples of 4096 bytes

Note: gzip is unused in practice

Note: uncompressed data is in NBT format, in chunk format(?)

### Format

* "sections" list tag with (up to 16) compound tags
* each section has 16x16x16 "Blocks", "Data", "SkyLight", "BlockLight" (chunk)
* each section has "Y" byte tag 0 bottom 15 top
* each section has optinoal "Add" tag, duplicate of "Data" (used to calculate blockid = (add << 8) + base)
* each chunk has 16x16 byte array "Biomes"
* new format y z x ((y * 16 + z) * 16 + x)
* ".mca" extension
* "Heightmap" tag uses NBT Int Array.

### NBT Format
Compound "": root tag
* Compound "Level": chunk data
  * Int "xPos": X position of chunk
  * Int "zPos": Z position of chunk
  * Long "LastUpdate": tick when chunk was last saved (0.0)
  * Byte "LightPopulated": unknown 1/0 (true/false) -> (0)
  * Byte "TerrainPopulated": have "special things" been added (ore, trees) 1/abs (true/false) -> (1)
  * Byte "V": 1 (likely chunk version tag)
  * Long "InhabitedTime": cumulative number of ticks players have been here (0)
  * Byte_Array "Biomes": 256 bytes, one per column, in what order?
  * Int_Array "HeightMap": 256 TAG_Int (all 0?)
  * List "Sections":
    * Byte "Y": index (not coordinate!) (0-15)
    * Byte_Array "Blocks": 4096 bytes of block IDs
    * Byte_Array "Add": optional, 2048 bytes of additional data (makes block values 12 bits long)
    * Byte_Array "Data": 2048 bytes of block data (4bits/block)
    * Byte-Array "BlockLight": 4bits/block (all 0?)
    * Byte-Array "SkyLight": 4bits/block (all 0?)
  * List "Entities": list of Compounds (do I care)
  * List "TileEntities": list of Compounds (ditto)
  * List "TileTicks": may not exist (so it won't)

## Blocks

Values shown for some items are not the default!

 * Wheat, Carrot, Potato, Pumpkin Stem, and Melon Stem are all full-grown

Some value range are not included because the server will handle it for us.

 * Redstone (wire, repeaters, comparators), Daylight Sensors, Farmland

Many other things are not supported because they are currently outside the scope of this project.  If someone else wants to add them, pull requests are accepted!

 * Torches and Redstone Torches
 * Piston and Piston Extension
 * Stairs
 * Beds
 * Signs
 * Doors
 * Rails
 * Ladders, Wall Signs, Furnaces, Chest (facing/attached)
 * Dispensers, Droppers, Hoppers (down, up, powered)
 * Levers (power, facing, wall/ground/ceiling)
 * Pressure Plates (pressed or not)
 * Buttons (set, direction)
 * Snow (depth)
 * Jukebox
 * Pumpkins and Jack o'Lanterns
 * Cake
 * Redstone Repeaters and Comparators
 * Trapdoors
 * Monster Eggs
 * Huge mushrooms
 * Vines
 * Fence Gates
 * Nether Wart
 * Brewing Stand
 * Cauldron
 * End Portal Block
 * Cocoa
 * Tripwire and Tripwire Hook
 * Heads

# Stuff to keep in mind

## Optimization

Check to see if a section is empty and just don't write it?

Skip writing any air blocks at first?

