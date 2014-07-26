// Minecraft world package.

package world

// Region files have names like "region/r.x.z.mcr" where x and z are
// "region coordinates".  These coordinates can be found this way:
// int regionX = chunkX >> 5
// int regionZ = chunkZ >> 5

// Region files begin with 8kiB header with what chunks are present,
// when last updated, and where found.
// Location in file of chunk at chunk coordinates (x, z) can be found at:
// byte offset 4 * (( x mod 32 ) + (z mod 32 ) * 32)
// (if foo mod 32 is negative, add 32)
// Timestamp of chunk is 4096 bytes later in file.
// bytes     description:
// 0-4095    locations (1024 entries)
// 4096-8191 timestamps (1024 entries)
// 8192+     chunks and unused space?

// chunk location
// byte description
// 0-2  offset (big-endian, in 4KiB sectors)
// 3    count (number of 4KiB sectors)
// Note: a chunk with an offset of 2 will begin right after timestamps table

// chunk timestamps
// byte description
// 0-3  timestamp (four-byte big-endian integer)

// chunk data
// byte description
// 0-3  length (in bytes)
// 4    compression type (1=gzip, 2=zlib)
// 5    compressed data (length-1 bytes)
// note: all chunks must be padded to multiples of 4096 bytes
//       gzip is unused in practice
//       uncompressed data is in NBT format, in chunk format

// chunk format
// - "sections" list tag with (up to 16) compound tags
// - each section has 16x16x16 "Blocks", "Data", "SkyLight", "BlockLight"
//   (chunk)
// - each section has "Y" byte tag 0 bottom 15 top
// - each section has optinoal "Add" tag, duplicate of "Data"
//   (used to calculate blockid = (add << 8) + base)
// - each chunk has 16x16 byte array "Biomes"
// - new format y z x ((y * 16 + z) * 16 + x)
// - ".mca" extension
// - "Heightmap" tag uses NBT Int Array.

// chunk NBT format structure:
// Compound "": root tag
// - Compound "Level": chunk data
//   - Int "xPos": X position of chunk
//   - Int "zPos": Z position of chunk
//   - Long "LastUpdate": tick when chunk was last saved (0.0)
//   - Byte "LightPopulated": unknown
//       1/0 (true/false) -> (0)
//   - Byte "TerrainPopulated": have "special things" been added (ore, trees)
//       1/abs (true/false) -> (1)
//   - Byte "V": 1 (likely chunk version tag)
//   - Long "InhabitedTime": cumulative number of ticks players have been here
//       (0)
//   - Byte_Array "Biomes": 256 bytes, one per column, in what order?
//   - Int_Array "HeightMap": 256 TAG_Int
//       (all 0?)
//   - List "Sections":
//     - Byte "Y": index (not coordinate!) (0-15)
//     - Byte_Array "Blocks": 4096 bytes of block IDs
//     - Byte_Array "Add": optional, 2048 bytes of additional data
//       (makes block values 12 bits long)
//     - Byte_Array "Data": 2048 bytes of block data (4bits/block)
//     - Byte-Array "BlockLight": 4bits/block (all 0?)
//     - Byte-Array "SkyLight": 4bits/block (all 0?)
//   - List "Entities": list of Compounds (do I care)
//   - List "TileEntities": list of Compounds (ditto)
//   - List "TileTicks": may not exist (so it won't)
