Development task list
---------------------

- [ ] Write NBT package
  - [ ] Implement whole spec now rather than wait
  - [ ] Use level.dat for test when possible
  - [ ] Double-check test coverage!
- [ ] Configure CI through drone.io
- [ ] Write world package
  - [ ] Write level.dat and test
  - [ ] Random access to blocks and data in world
  - [ ] Double-check test coverage!
- [ ] Write map generation package including biome map
  - [ ] Import landcover data into map
  - [ ] Generate crust array
  - [ ] Biome map: ocean/plains!
  - [ ] Generate bathy array
  - [ ] Biome map: deep ocean
  - [ ] Import elevation data into map
  - [ ] Biome map: Hills and Extreme hills
- [ ] Additional biome map features
  - [ ] Forests
  - [ ] Deserts
  - [ ] Developed land (buildings?  stone?)
  - [ ] Croplands
  - [ ] Beaches
  - [ ] What else?
- [ ] Translate biome and other maps into Minecraft world
  - [ ] One goroutine per chunk?  column?  region?  world?
- [ ] Build Docker container for project
- [ ] Integration tests -- known input, known outputs

Future
------
- [ ] Reuse world package as terrain generator for other servers
- [ ] Support for villages
- [ ] Support for rivers
- [ ] Use bathymetric data instead of guesses
- [ ] Use transportation data (railway network?  road?)
