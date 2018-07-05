package main

// world dimension of the grid (in world-space coordinates,
// the dimension of the grid square that is the screen)
const GRID_WORLD_DIMENSION = 1024

// the grid is made of GRID_CELL_DIMENSION x GRID_CELL_DIMENSION cells
const GRID_CELL_DIMENSION = 12

// visual constants
const WINDOW_WIDTH = 640
const WINDOW_HEIGHT = 640
const FONTSZ = 16
const FPS = 60

// dimensions of grid cells in pixels
const GRIDCELL_PX_W = WINDOW_WIDTH / GRID_CELL_DIMENSION
const GRIDCELL_PX_H = WINDOW_HEIGHT / GRID_CELL_DIMENSION

// dimensions of grid cells in world space
const GRIDCELL_WORLD_W = GRID_WORLD_DIMENSION / GRID_CELL_DIMENSION
const GRIDCELL_WORLD_H = GRID_WORLD_DIMENSION / GRID_CELL_DIMENSION

// density with which to populate obstacles
const DENSITY = 0.05
