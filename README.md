[![Build Status](https://drone.io/github.com/mathuin/terroir/status.png)](https://drone.io/github.com/mathuin/terroir/latest)

# Welcome to Terroir!

Terroir uses real data to generate Minecraft worlds which somewhat resemble actual places in the real world.

## Installation

**TODO:** Install Docker, and fetch the pre-built image.

## Requirements

Terroir requires a certain basic set of data in order to generate worlds.

* Raw mapping data
  * Elevation (1/3 arc-second is fine)
  * Landcover (NLCD 2011 is fine)
* Map parameters
  * Projection
  * Landcover translation table
* Local parameters
  * Latlong boundaries
  * Scaling values

**TODO:** Use whatever tools you need to use to retrieve your mapping data.  Example of National Map viewer goes here.

**TODO:** _Any_ example from another country would be very nice!

## Preparation

**TODO:** Run the batch file generator script against your settings.

## Execution

**TODO:** Run the pre-built Docker container with the following parameters.

**TODO:** The output is a ready to use Minecraft world!

## Future work

I am accepting pull requests for bugfixes and features!

I would like to see (and hope to add) additional support for the following:

* Non-US regions (specifically EU)
* Rivers
* Villages

That being said, it is highly unlikely that I will personally add the following features:

* Iconic buildings (bring your own Eiffel Tower!)
* Other planets or moons

I will still accept pull requests for those features, they just aren't interesting enough to me for me to add them!

# License

This software is released under the [MIT license](http://opensource.org/licenses/MIT) as documented [here](../blob/master/LICENSE.md).
