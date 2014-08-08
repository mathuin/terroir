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
- [ ] Write map generation package including biome map
  - [ ] Import landcover data into map
  - [ ] Generate crust array
  - [ ] Biome map: ocean/plains!
  - [ ] Generate bathy array
  - [ ] Biome map: deep ocean
  - [ ] Import elevation data into map
  - [ ] Biome map: Hills and Extreme hills
- [ ] Translate biome and other maps into Minecraft world
  - [ ] One goroutine per chunk?  column?  region?  world?
- [ ] Build Docker container for project
- [ ] Integration tests -- known input, known outputs

Future
------
- [ ] Additional biome map features
  - [ ] Forests
  - [ ] Deserts
  - [ ] Buildings on developed lands
  - [ ] Croplands
  - [ ] Beaches
  - [ ] Rivers
  - [ ] What else?
- [ ] Reuse world package as terrain generator for other servers
- [ ] Use bathymetric data instead of guesses
- [ ] Use transportation data (railway network?  road?)
- [ ] Villages in developed lands
