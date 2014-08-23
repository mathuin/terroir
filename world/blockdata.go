package world

func init() {
	for _, bd := range blockData {
		blockNames[bd.name] = MakeBlock(bd.block, bd.data)
	}
}

var blockData = []struct {
	block int
	data  int
	name  string
}{
	{0, 0, "Air"},
	{1, 0, "Stone"},
	{1, 1, "Granite"},
	{1, 2, "Polished Granite"},
	{1, 3, "Diorite"},
	{1, 4, "Polished Diorite"},
	{1, 5, "Andesite"},
	{1, 6, "Polished Andesite"},
	{2, 0, "Grass Block"},
	{3, 0, "Dirt"},
	{3, 1, "Coarse Dirt"},
	{3, 2, "Podzol"},
	{4, 0, "Cobblestone"},
	{5, 0, "Oak Wood Planks"},
	{5, 1, "Spruce Wood Planks"},
	{5, 2, "Birch Wood Planks"},
	{5, 3, "Jungle Wood Planks"},
	{5, 4, "Acacia Wood Planks"},
	{5, 5, "Dark Oak Wood Planks"},
	{6, 0, "Oak Sapling"},
	{6, 1, "Spruce Sapling"},
	{6, 2, "Birch Sapling"},
	{6, 3, "Jungle Sapling"},
	{6, 4, "Acacia Sapling"},
	{6, 5, "Dark Oak Sapling"},
	{7, 0, "Bedrock"},
	{8, 0, "Flowing Water"},
	{9, 0, "Water"},
	{10, 0, "Flowing Lava"},
	{11, 0, "Lava"},
	{12, 0, "Sand"},
	{12, 1, "Red Sand"},
	{13, 0, "Gravel"},
	{14, 0, "Gold Ore"},
	{15, 0, "Iron Ore "},
	{16, 0, "Coal Ore"},
	{17, 0, "Oak Wood (Vertical)"},
	{17, 1, "Spruce Wood (Vertical)"},
	{17, 2, "Birch Wood (Vertical)"},
	{17, 3, "Jungle Wood (Vertical)"},
	{17, 4, "Oak Wood (East/West)"},
	{17, 5, "Spruce Wood (East/West)"},
	{17, 6, "Birch Wood (East/West)"},
	{17, 7, "Jungle Wood (East/West)"},
	{17, 8, "Oak Wood (North/South)"},
	{17, 9, "Spruce Wood (North/South)"},
	{17, 10, "Birch Wood (North/South)"},
	{17, 11, "Jungle Wood (North/South)"},
	{17, 12, "Oak Wood (Bark)"},
	{17, 13, "Spruce Wood (Bark)"},
	{17, 14, "Birch Wood (Bark)"},
	{17, 15, "Jungle Wood (Bark)"},
	{18, 0, "Oak Leaves"},
	{18, 1, "Spruce Leaves"},
	{18, 2, "Birch Leaves"},
	{18, 3, "Jungle Leaves"},
	{18, 4, "Oak Leaves (No Decay)"},
	{18, 5, "Spruce Leaves (No Decay)"},
	{18, 6, "Birch Leaves (No Decay)"},
	{18, 7, "Jungle Leaves (No Decay)"},
	{18, 8, "Oak Leaves (Check Decay)"},
	{18, 9, "Spruce Leaves (Check Decay)"},
	{18, 10, "Birch Leaves (Check Decay)"},
	{18, 11, "Jungle Leaves (Check Decay)"},
	{18, 12, "Oak Leaves (No Decay And Check Decay)"},
	{18, 13, "Spruce Leaves (No Decay And Check Decay)"},
	{18, 14, "Birch Leaves (No Decay And Check Decay)"},
	{18, 15, "Jungle Leaves (No Decay And Check decay)"},
	{19, 0, "Sponge"},
	{19, 1, "Wet Sponge"},
	{20, 0, "Glass"},
	{21, 0, "Lapis Lazuli Ore"},
	{22, 0, "Lapis Lazuli Block"},
	{23, 0, "Dispenser"},
	{24, 0, "Sandstone"},
	{24, 1, "Chiseled Sandstone"},
	{24, 2, "Smooth Sandstone"},
	{25, 0, "Note Block"},
	{26, 0, "Bed"},
	{27, 0, "Powered Rail"},
	{28, 0, "Detector Rail"},
	{29, 0, "Sticky Piston"},
	{30, 0, "Cobweb"},
	{31, 0, "Shrub"},
	{31, 1, "Grass"},
	{31, 2, "Fern"},
	{31, 3, "Grass (Biome)"},
	{32, 0, "Dead Bush"},
	{33, 0, "Piston"},
	{34, 0, "Piston Extension"},
	{35, 0, "White Wool"},
	{35, 1, "Orange Wool"},
	{35, 2, "Magenta Wool"},
	{35, 3, "Light Blue Wool"},
	{35, 4, "Yellow Wool"},
	{35, 5, "Lime Wool"},
	{35, 6, "Pink Wool"},
	{35, 7, "Gray Wool"},
	{35, 8, "Light Gray Wool"},
	{35, 9, "Cyan Wool"},
	{35, 10, "Purple Wool"},
	{35, 11, "Blue Wool"},
	{35, 12, "Brown Wool"},
	{35, 13, "Green Wool"},
	{35, 14, "Red Wool"},
	{35, 15, "Black Wool"},
	{36, 0, "Block moved by Piston"},
	{37, 0, "Dandelion"},
	{38, 0, "Poppy"},
	{38, 1, "Blue Orchid"},
	{38, 2, "Allium"},
	{38, 3, "Azure Bluet"},
	{38, 4, "Red Tulip"},
	{38, 5, "Orange Tulip"},
	{38, 6, "White Tulip"},
	{38, 7, "Pink Tulip"},
	{38, 8, "Oxeye Daisy"},
	{39, 0, "Brown Mushroom"},
	{40, 0, "Red Mushroom"},
	{41, 0, "Block of Gold"},
	{42, 0, "Block of Iron"},
	{43, 0, "Double Stone Slab"},
	{43, 1, "Double Sandstone Slab"},
	{43, 2, "Double (Stone) Wooden Slab"},
	{43, 3, "Double Cobblestone Slab"},
	{43, 4, "Double Bricks Slab"},
	{43, 5, "Double Stone Brick Slab"},
	{43, 6, "Double Nether Brick Slab"},
	{43, 7, "Double Quartz Slab"},
	{43, 8, "Full Stone Slab"},
	{43, 9, "Full Sandstone Slab"},
	{43, 10, "Tile Quartz Slab"},
	{44, 0, "Stone Slab"},
	{44, 1, "Sandstone Slab"},
	{44, 2, "(Stone) Wooden Slab"},
	{44, 3, "Cobblestone Slab"},
	{44, 4, "Bricks Slab"},
	{44, 5, "Stone Brick Slab"},
	{44, 6, "Nether Brick Slab"},
	{44, 7, "Quartz Slab"},
	{44, 8, "Upper Stone Slab"},
	{44, 9, "Upper Sandstone Slab"},
	{44, 10, "Upper (Stone) Wooden Slab"},
	{44, 11, "Upper Cobblestone Slab"},
	{44, 12, "Upper Bricks Slab"},
	{44, 13, "Upper Stone Brick Slab"},
	{44, 14, "Upper Nether Brick Slab"},
	{44, 15, "Upper Quartz Slab"},
	{45, 0, "Bricks"},
	{46, 0, "TNT"},
	{47, 0, "Bookshelf "},
	{48, 0, "Moss Stone"},
	{49, 0, "Obsidian"},
	{50, 0, "Torch"},
	{51, 0, "Fire"},
	{52, 0, "Monster Spawner"},
	{53, 0, "Oak Wood Stairs"},
	{54, 0, "Chest"},
	{55, 0, "Redstone Wire"},
	{56, 0, "Diamond Ore"},
	{57, 0, "Block of Diamond"},
	{58, 0, "Crafting Table"},
	{59, 7, "Wheat"},
	{60, 0, "Farmland"},
	{61, 0, "Furnace"},
	{62, 0, "Burning Furnace"},
	{63, 0, "Standing Sign "},
	{64, 0, "Wooden Door"},
	{65, 0, "Ladder"},
	{66, 0, "Rail"},
	{67, 0, "Cobblestone Stairs"},
	{68, 0, "Wall Sign"},
	{69, 0, "Lever"},
	{70, 0, "Stone Pressure Plate"},
	{71, 0, "Iron Door"},
	{72, 0, "Wooden Pressure Plate"},
	{73, 0, "Redstone Ore"},
	{74, 0, "Glowing Redstone Ore"},
	{75, 0, "Redstone Torch (inactive)"},
	{76, 0, "Redstone Torch (active)"},
	{77, 0, "Stone Button"},
	{78, 0, "Snow"},
	{79, 0, "Ice "},
	{80, 0, "Snow"},
	{81, 0, "Cactus"},
	{82, 0, "Clay"},
	{83, 0, "Sugar Cane"},
	{84, 0, "Jukebox"},
	{85, 0, "Fence"},
	{86, 0, "Pumpkin"},
	{87, 0, "Netherrack"},
	{88, 0, "Soul Sand"},
	{89, 0, "Glowstone"},
	{90, 0, "Nether Portal"},
	{91, 0, "Jack o'Lantern"},
	{92, 0, "Cake"},
	{93, 0, "Redstone Repeater (Inactive)"},
	{94, 0, "Redstone Repeater (Active)"},
	{95, 0, "White Stained Glass"},
	{95, 1, "Orange Stained Glass"},
	{95, 2, "Magenta Stained Glass"},
	{95, 3, "Light Blue Stained Glass"},
	{95, 4, "Yellow Stained Glass"},
	{95, 5, "Lime Stained Glass"},
	{95, 6, "Pink Stained Glass"},
	{95, 7, "Gray Stained Glass"},
	{95, 8, "Light Gray Stained Glass"},
	{95, 9, "Cyan Stained Glass"},
	{95, 10, "Purple Stained Glass"},
	{95, 11, "Blue Stained Glass"},
	{95, 12, "Brown Stained Glass"},
	{95, 13, "Green Stained Glass"},
	{95, 14, "Red Stained Glass"},
	{95, 15, "Black Stained Glass"},
	{96, 0, "Trapdoor"},
	{97, 0, "Monster Egg"},
	{98, 0, "Stone Bricks"},
	{98, 1, "Mossy Stone Bricks"},
	{98, 2, "Cracked Stone Bricks"},
	{98, 3, "Chiseled Stone Bricks"},
	{99, 0, "Huge Brown Mushroom"},
	{100, 0, "Huge Red Mushroom"},
	{101, 0, "Iron Bars"},
	{102, 0, "Glass Pane"},
	{103, 0, "Melon"},
	{104, 7, "Pumpkin Stem"},
	{105, 7, "Melon Stem"},
	{106, 0, "Vines"},
	{107, 0, "Fence Gate"},
	{108, 0, "Brick Stairs"},
	{109, 0, "Stone Brick Stairs"},
	{110, 0, "Mycelium"},
	{111, 0, "Lily Pad "},
	{112, 0, "Nether Brick"},
	{113, 0, "Nether Brick Fence"},
	{114, 0, "Nether Brick Stairs"},
	{115, 0, "Nether Wart"},
	{116, 0, "Enchantment Table"},
	{117, 0, "Brewing Stand"},
	{118, 0, "Cauldron"},
	{119, 0, "End Portal"},
	{120, 0, "End Portal Block"},
	{121, 0, "End Stone"},
	{122, 0, "Dragon Egg"},
	{123, 0, "Redstone Lamp (inactive)"},
	{124, 0, "Redstone Lamp (active)"},
	{125, 0, "Double Oak Wood Slab"},
	{125, 1, "Double Spruce Wood Slab"},
	{125, 2, "Double Birch Wood Slab"},
	{125, 3, "Double Jungle Wood Slab"},
	{125, 4, "Double Acacia Wood Slab"},
	{125, 5, "Double Dark Oak Wood Slab"},
	{126, 0, "Oak Wood Slab"},
	{126, 1, "Spruce Wood Slab"},
	{126, 2, "Birch Wood Slab"},
	{126, 3, "Jungle Wood Slab"},
	{126, 4, "Acacia Wood Slab"},
	{126, 5, "Dark Oak Wood Slab"},
	{126, 8, "Upper Oak Wood Slab"},
	{126, 9, "Upper Spruce Wood Slab"},
	{126, 10, "Upper Birch Wood Slab"},
	{126, 11, "Upper Jungle Wood Slab"},
	{126, 12, "Upper Acacia Wood Slab"},
	{126, 13, "Upper Dark Oak Wood Slab"},
	{127, 0, "Cocoa "},
	{128, 0, "Sandstone Stairs"},
	{129, 0, "Emerald Ore"},
	{130, 0, "Ender Chest"},
	{131, 0, "Tripwire Hook"},
	{132, 0, "Tripwire"},
	{133, 0, "Block of Emerald"},
	{134, 0, "Spruce Wood Stairs"},
	{135, 0, "Birch Wood Stairs"},
	{136, 0, "Jungle Wood Stairs"},
	{137, 0, "Command Block"},
	{138, 0, "Beacon"},
	{139, 0, "Cobblestone Wall"},
	{139, 1, "Mossy Cobblestone Wall"},
	{140, 0, "Flower Pot"},
	{141, 7, "Carrot"},
	{142, 7, "Potato"},
	{143, 0, "Wooden Button "},
	{144, 0, "Mob Head"},
	{145, 0, "Anvil (North/South)"},
	{145, 1, "Anvil (East/West)"},
	{145, 2, "Anvil (South/North)"},
	{145, 3, "Anvil (West/East)"},
	{145, 4, "Slightly Damaged Anvil (North/South)"},
	{145, 5, "Slightly Damaged Anvil (East/West)"},
	{145, 6, "Slightly Damaged Anvil (South/North)"},
	{145, 7, "Slightly Damaged Anvil (West/East)"},
	{145, 8, "Very Damaged Anvil (North/South)"},
	{145, 9, "Very Damaged Anvil (East/West)"},
	{145, 10, "Very Damaged Anvil (South/North)"},
	{145, 11, "Very Damaged Anvil (West/East)"},
	{146, 0, "Trapped Chest"},
	{147, 0, "Light Weighted Pressure Plate"},
	{148, 0, "Heavy Weighted Pressure Plate"},
	{149, 0, "Redstone Comparator (Unpowered)"},
	{150, 0, "Redstone Comparator (Powered)"},
	{151, 0, "Daylight Sensor"},
	{152, 0, "Block of Redstone"},
	{153, 0, "Nether Quartz Ore"},
	{154, 0, "Hopper"},
	{155, 0, "Block of Quartz"},
	{155, 1, "Chiseled Block of Quartz"},
	{155, 4, "Pillar Quartz Block (Vertical)"},
	{155, 4, "Pillar Quartz Block (North-South)"},
	{155, 4, "Pillar Quartz Block (East-West)"},
	{156, 0, "Quartz Stairs"},
	{157, 0, "Activator Rail"},
	{158, 0, "Dropper"},
	{159, 0, "White Stained Clay"},
	{159, 1, "Orange Stained Clay"},
	{159, 2, "Magenta Stained Clay"},
	{159, 3, "Light Blue Stained Clay"},
	{159, 4, "Yellow Stained Clay"},
	{159, 5, "Lime Stained Clay"},
	{159, 6, "Pink Stained Clay"},
	{159, 7, "Gray Stained Clay"},
	{159, 8, "Light Gray Stained Clay"},
	{159, 9, "Cyan Stained Clay"},
	{159, 10, "Purple Stained Clay"},
	{159, 11, "Blue Stained Clay"},
	{159, 12, "Brown Stained Clay"},
	{159, 13, "Green Stained Clay"},
	{159, 14, "Red Stained Clay"},
	{159, 15, "Black Stained Clay"},
	{160, 0, "White Stained Glass Pane"},
	{160, 1, "Orange Stained Glass Pane"},
	{160, 2, "Magenta Stained Glass Pane"},
	{160, 3, "Light Blue Stained Glass Pane"},
	{160, 4, "Yellow Stained Glass Pane"},
	{160, 5, "Lime Stained Glass Pane"},
	{160, 6, "Pink Stained Glass Pane"},
	{160, 7, "Gray Stained Glass Pane"},
	{160, 8, "Light Gray Stained Glass Pane"},
	{160, 9, "Cyan Stained Glass Pane"},
	{160, 10, "Purple Stained Glass Pane"},
	{160, 11, "Blue Stained Glass Pane"},
	{160, 12, "Brown Stained Glass Pane"},
	{160, 13, "Green Stained Glass Pane"},
	{160, 14, "Red Stained Glass Pane"},
	{160, 15, "Black Stained Glass Pane"},
	{161, 0, "Acacia Leaves"},
	{161, 1, "Dark Oak Leaves"},
	{161, 4, "Acacia Leaves (No Decay)"},
	{161, 5, "Dark Oak Leaves (No Decay)"},
	{161, 8, "Acacia Leaves (Check Decay)"},
	{161, 9, "Dark Oak Leaves (Check Decay)"},
	{161, 12, "Acacia Leaves (No Decay And Check Decay)"},
	{161, 13, "Dark Oak Leaves (No Decay And Check Decay)"},
	{162, 0, "Acacia Wood (Vertical)"},
	{162, 1, "Dark Oak Wood (Vertical)"},
	{162, 4, "Acacia Wood (East/West)"},
	{162, 5, "Dark Oak Wood (East/West)"},
	{162, 8, "Acacia Wood (North/South)"},
	{162, 9, "Dark Oak Wood (North/South)"},
	{162, 12, "Acacia Wood (Bark)"},
	{162, 13, "Dark Oak Wood (Bark)"},
	{163, 0, "Acacia Wood Stairs"},
	{164, 0, "Dark Oak Wood Stairs"},
	{165, 0, "Slime Block"},
	{166, 0, "Barrier"},
	{167, 0, "Iron Trapdoor"},
	{168, 0, "Prismarine"},
	{168, 1, "Prismarine Bricks"},
	{168, 2, "Dark Prismarine"},
	{169, 0, "Sea Lantern"},
	{170, 0, "Hay Block"},
	{171, 0, "White Carpet"},
	{171, 1, "Orange Carpet"},
	{171, 2, "Magenta Carpet"},
	{171, 3, "Light Blue Carpet"},
	{171, 4, "Yellow Carpet"},
	{171, 5, "Lime Carpet"},
	{171, 6, "Pink Carpet"},
	{171, 7, "Gray Carpet"},
	{171, 8, "Light Gray Carpet"},
	{171, 9, "Cyan Carpet"},
	{171, 10, "Purple Carpet"},
	{171, 11, "Blue Carpet"},
	{171, 12, "Brown Carpet"},
	{171, 13, "Green Carpet"},
	{171, 14, "Red Carpet"},
	{171, 15, "Black Carpet"},
	{173, 0, "Block of Coal"},
	{174, 0, "Packed Ice"},
	{175, 0, "Sunflower"},
	{175, 1, "Lilac"},
	{175, 2, "Double Tallgrass"},
	{175, 3, "Large Fern"},
	{175, 4, "Rose Bush"},
	{175, 5, "Peony"},
	{175, 8, "Sunflower Top"},
	{175, 9, "Lilac Top"},
	{175, 10, "Double Tallgrass Top"},
	{175, 11, "Large Fern Top"},
	{175, 12, "Rose Bush Top"},
	{175, 13, "Peony Top"},
	{176, 0, "Standing Banner"},
	{177, 0, "Wall Banner"},
	{178, 0, "Inverted Daylight Sensor"},
	{179, 0, "Red Sandstone"},
	{179, 1, "Chiseled Red Sandstone"},
	{179, 2, "Smooth Red Sandstone"},
	{181, 0, "Double Red Sandstone Slab"},
	{181, 8, "Full Red Sandstone Slab"},
	{182, 0, "Red Sandstone Slab"},
	{182, 8, "Upper Red Sandstone Slab"},
}
