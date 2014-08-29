Development task list
---------------------

- [x] Write NBT package
  - [x] Implement whole spec now rather than wait
  - [x] Use level.dat for test when possible
  - [x] Double-check test coverage!
- [x] Configure CI through drone.io
- [x] Write world package
  - [x] Write level.dat and test
  - [x] Random access to blocks and data in world
  - [x] Double-check test coverage!
- [x] Write map generation package including biome map
  - [x] Import landcover data into map
  - [x] Generate crust array
  - [x] Biome map: ocean/plains!
  - [x] Generate bathy array
  - [x] Biome map: deep ocean
  - [x] Import elevation data into map
  - [x] Biome map: Hills and Extreme hills
- [x] Translate biome and other maps into Minecraft world
  - [x] One goroutine per chunk?  column?  region?  world?
- [ ] Build Docker container for project
- [ ] Integration tests -- known input, known outputs

Future
------
- [ ] Additional biome map features
  - [x] Forests
  - [x] Deserts
  - [ ] Buildings on developed lands
  - [ ] Croplands
  - [ ] Beaches
  - [ ] Rivers
  - [ ] What else?
- [ ] Reuse world package as terrain generator for other servers
- [ ] Use bathymetric data instead of guesses
- [ ] Use canopy and impervious surface data instead of guesses
- [ ] Use transportation data (railway network?  road?)
- [ ] Villages in developed lands
