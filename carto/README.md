# Carto notes

The purpose of this package is to take the map data supplied by the user and construct a multi-band GeoTIFF which can be used by the build package to construct a Minecraft world.

## Problems

Why is the map reflected across y=-x or the NW/SE axis?  What's north in qgis is west in Minecraft, and what's north in Minecraft is west in qgis.  It's possible this is a rotation/skew sort of thing.  Check region.py for how I fixed it there.

In level=30 on the full Block Island render, nodata is not being handled correctly.  Band 1 is 0 when it should be 11.  Band 3 is 0 when it should be .. maxdepth, but that should fix itself at the time

Parallelize (correctly!) the code that checks what points are in the shape.

Change the Out type to Column, with the struct to include the XZ value, the biome value, and an array from 0 to wherever of the blocks.  Change the write routine to handle this correctly.

Be *slightly* less lazy with how we're rendering the columns for demonstration purposes.  Do the clever with forests, for example, and ocean/deep-ocean.

Knob all the things, and return the make-region stuff to its previous state.

Finally: remove the unused code because it doesn't have much to teach us.


