# gocube

A tiny Go Rubik's Cube solver and visualizer.

## Run

```
go run .
```

## Summary

Main function currently: scrambles a cube, prints the cube, then solves it with BFS.

## Limitaion

Maximum solve depth for an optimal solve is 5 turns, above that is too high for 16GB of ram on windows.

## Plans

1. Implement a way for the user to enter a cube state.

2. Create 3D visualization using raylib or WebGPU.

3. Port code to mobile device.

4. Add scanning ability for both web and mobile app

## Problems

Currently the LUT for the corners is calculated incorrectly as it throws away shorter routes to nodes and overestimates the distance in the heuristic, which no longer guarantees optimal solution.

At the moment an alternate RCube state is being created to ideally decrease the cost of finding the rank of some state. The current implementation does not properly calculate the new orientaition rank, as permutation affects the orientation transitions.