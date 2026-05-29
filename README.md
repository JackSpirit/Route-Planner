OpenStreetMap Route Planner & Reachability Engine

A high-performance, lightweight graph-based routing engine built from scratch in Go, paired with an interactive Leaflet.js map frontend. The project parses raw OpenStreetMap (OSM) data into an in-memory graph structure, allowing users to calculate optimal point-to-point routes and visualize reachable road networks within dynamic time limits (Isochrones).
Features

    Custom OSM Parser: Processes node and way data structures into a spatial graph layout.

    Dual Routing Profiles:

        Shortest Distance: Calculates paths optimizing strictly for minimum physical mileage.

        Fastest Time: Weighs edges using road classification speed limits (motorways, trunk roads, residential, etc.) to optimize for travel duration.

    High-Performance Dijkstra Implementation: Implements a custom minimum binary heap wrapper (container/heap) supporting fast O((E + V) log V) node updates and state management.

    Dynamic Network Isochrones: Computes a full reachability analysis from any point. Instead of inaccurate geometric bounding shapes, it paints the exact tree branch skeleton of navigable roads reachable under your chosen time threshold.

    Interactive Leaflet Map Control UI: Provides an integrated HUD widget for seamless routing profile switching and live travel-time adjustments.

Project Structure

    app.js: Frontend logic (Leaflet initialization, event handlers, API fetch)

    index.html: Map viewport container and DOM structure

    main.go: Entry point (initializes data parsing and launches components)

    graph.go: Core graph generation, distance formulas, and edge metrics

    queue.go: Custom Min-Priority Queue implementation tracking path item weights

    router.go: Core point-to-point Dijkstra shortest-path engine

    isochrone.go: Multi-destination exploratory BFS/Dijkstra reachability algorithm

    server.go: HTTP Server providing REST API endpoints (/route and /isochrone)


    
## Backend API Specifications

The Go backend spins up an HTTP server on `http://localhost:8080`.

### 1. Point-to-Point Routing

* **Endpoint:** `/route`
* **Method:** `GET`
* **Query Parameters:**
* `start_lat` / `start_lon`: Floating point coordinates for origin node.
* `end_lat` / `end_lon`: Floating point coordinates for destination node.
* `mode`: `distance` or `time`.


* **Sample Response:**
```json
{
  "path": [
    {"lat": 39.1512, "lon": -75.5234},
    {"lat": 39.1545, "lon": -75.5198}
  ],
  "distance_km": 4.82,
  "node_count": 42
}

```



### 2. Network Isochrone (Reachability Tree)

* **Endpoint:** `/isochrone`
* **Method:** `GET`
* **Query Parameters:**
* `lat` / `lon`: Central origin coordinate.
* `minutes`: Maximum travel time cutoff (e.g., `10`).


* **Sample Response:**
```json
{
  "segments": [
    [[39.1512, -75.5234], [39.1523, -75.5211]],
    [[39.1523, -75.5211], [39.1541, -75.5195]]
  ]
}

```



---

## Getting Started

### Prerequisites

* Go (version 1.18 or higher recommended)
* A modern web browser with internet connectivity (to pull map tiles and Leaflet assets)

### Installation & Run

1. **Clone or navigate to your project directory:**
```bash
cd Route-Planner

```


2. **Compile and execute the Go backend package:**
```bash
go run *.go

```


*You should see the confirmation terminal log:* `Web server starting on http://localhost:8080..`
3. **Launch the frontend:**
Open `index.html` directly in your browser or run it through a local web server tool (e.g., Live Server extension in VS Code).

---

## How to Use the Map

### Finding Standard Routes

1. Set your **Routing Profile** to *Shortest Distance* or *Fastest Time* using the top-right map HUD.
2. **Left-click** anywhere on the map to set your **Start point** (Green marker).
3. **Left-click** a second location to set your **End point** (Red marker).
4. The system calculates the route instantly and draws a solid dark blue navigation polyline tracking your chosen metrics.
5. Clicking a third time will clear the map state for a new path.

### Visualizing Reachable Roads (Isochrones)

1. Use the **Isochrone Limit** input box in the top-right control container to choose a maximum time budget in minutes (e.g., `10` or `15`).
2. Hold down the **Shift key** on your keyboard.
3. **Left-click** your target starting location on the map.
4. The engine executes an exploratory frontier search along your road matrix and dynamically returns a glowing **blue network tree** visualizing every street segment passable within that absolute time frame.

```

```
