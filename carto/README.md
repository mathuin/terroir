# Carto notes

The purpose of this package is to take the map data supplied by the
user and construct a multi-band GeoTIFF which can be used by the build
package to construct a Minecraft world.

# Issues

## Coordinates

Map coordinates are different than Minecraft coordinates.  Details:

+---------+-----+-----+----+----+
|         |North|South|East|West|
+---------+-----+-----+----+----+
|Qgis     | +Y  | -Y  | +X | -X |
|Minecraft| -Z  | +Z  | +X | -X | 
+---------+-----+-----+----+----+

Right now, processPoints will correct by dividing the individual
coordinates by their transform values.  

## Performance

I have parallelized the code that checks what points are in the shape, but that doesn't really help.  The next step is OpenCL.

https://github.com/pseudomind/go-opencl actually works pretty well on nala now that I nuked the drivers. :-P

Inputs would be:  the whole array, or the generated list of polygons and the dimensions of the array.

Outputs would be: lists of all points in each polygon.

Rough idea on how to do it:

- for each row of envelope
   - set edge state to "left" ("right", "on" are other options)
     for all edges (outer or inner) which this row might cross
	 (i.e., whose maxX is >= row and minX is <= row)
   - set polygon to "0"
   - for each pixel of row
     - for each edge this row might cross (in "clockwise" order?)
	   - calculate sides for -0.5/+0.5
	     - if left/right:
		   set polygon to this polygon's counter
		 - if !on/on:
		   set polygon to this polygon's counter
		 - if on/!on
		   set polygon to previous polygon's counter

http://alienryderflex.com/polygon/ looks like someone's done something very much like it already.  Now to figure out the most efficient way to do this in OpenCL, handling the outer-inner ring thing as well.

Maybe sort by largest envelope and do those points first, skipping over the ones that are already done?

Someone else must have done this already!

### New idea

Inputs: bounding box, number of edges, array of edges, array for output (maximum size of bounding box)
Outputs: length used, array for output

## When it all works right...

Remove the unused code because it doesn't have much to teach us.


