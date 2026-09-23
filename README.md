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

1. Switch to Iterative Deepening A* with heuristic LUT.

2. Implement a way for the user to enter a cube state.

3. Create 3D visualization using raylib or WebGPU.

4. Port code to mobile device.

5. Add scanning ability for both web and mobile app