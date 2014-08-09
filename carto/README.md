# Map notes

## Get elevation and landcover data into arrays

For now, use the VRTs from the giant download of all that data for
production.  Once it's working, extract very very small chunks of both
datasets and save them in the repository along with a test routine.

## Build a biome map from that data

The initial map will be oceans/plains.  If landcover says it's water,
biome says it's ocean.  More sophisticated once everything's working.

## Build a three-dimensional array of blocks based on that data

For each XZ pair, construct a column that expresses the data in the
arrays.  Should be fun!

## Write those blocks to a Minecraft level

Might have to make chunk-level versions of SetBlock(), not sure yet.

# Stuff to keep in mind

## Optimizations

Maybe do SetBlock() on a section level?

Check to see if a section is empty and just don't write it?

Skip writing any air blocks at first?

