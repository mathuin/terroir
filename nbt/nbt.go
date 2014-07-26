// Named Binary Tag package

// The NBT format is documented in many places including here:
// http://minecraft.gamepedia.com/NBT_Format

package nbt

// version 19133, no need to support anything older

// tag format:
// byte 0: id
// byte 1-2: name len (TAG_End has no name so just one byte)
// byte 3-n: name UTF-8 (*may* contain spaces, but don't)
// byte (n+1)-m: payload

//  0 TAG_End
//    No payload, no name, just 0
//  1 TAG_Byte
//    1 byte / 8 bits, signed (also used for booleans sometimes)
//  2 TAG_Short
//    2 bytes / 16 bits, signed, big endian
//  3 TAG_Int
//    4 bytes / 32 bits, signed, big endian
//  4 TAG_Long
//    8 bytes / 64 bits, signed, big endian
//  5 TAG_Float
//    4 bytes / 32 bits, signed, big endian, IEEE 754-2008, binary32
//  6 TAG_Double
//    8 bytes / 64 bits, signed, big endian, IEEE 754-2008, binary64
//  7 TAG_Byte_Array
//    TAG_Int's payload "size", then size * TAG_Byte payloads
//  8 TAG_String
//    TAG_Short's payload "length", then length * UTF-8 code points
//  9 TAG_List
//    TAG_Byte's payload "tagId", then TAG_Int's payload "size",
//    then "size" tag's payloads, all of type tagId
// 10 TAG_Compound
//    Fully-formed tags followed by TAG_End.
// 11 TAG_Int_Array
//    TAG_Int's payload "size", then size * TAG_Int payloads

// Maximum nesting limit for List and Compound tags is 512.

// Proper NBT files are gzipped TAG_Compound tags with name and tag ID.

// First: figure out how to read them, using level.dat and friends.
// Then: figure out how to write them, comparing them to what was read.
// Next: write tests that match the output of individual tags.
// Next: write test that handle list and compound tags.
// Finally: finish dealing with nested files, match to level.dat
