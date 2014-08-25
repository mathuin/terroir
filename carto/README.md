# Carto notes

The purpose of this package is to take the map data supplied by the user and construct a multi-band GeoTIFF which can be used by the build package to construct a Minecraft world.

## Problems

The IDT code does not appear to be doing the right thing with the
landcover on the shoreline.

Why is the map reflected across y=-x or the NW/SE axis?  What's north in qgis is west in Minecraft, and what's north in Minecraft is west in qgis.  It's possible this is a rotation/skew sort of thing.  Check region.py for how I fixed it there.  

JMT:  I looked at the map in Qgis and the dynmap and this is what I saw.

+---------+-----+-----+----+----+
|         |North|South|East|West|
+---------+-----+-----+----+----+
|Qgis     | +Y  | -Y  | +X | -X |
|Minecraft| -Z  | +Z  | +X | -X | 
+---------+-----+-----+----+----+

Parallelize (correctly!) the code that checks what points are in the shape.

Knob all the things, and return the make-region stuff to its previous state.

Finally: remove the unused code because it doesn't have much to teach us.


