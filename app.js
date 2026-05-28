const map = L.map('map').setView([39.15, -75.52], 9);

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '© OpenStreetMap contributors'
}).addTo(map);

let startMarker = null;
let endMarker = null;
let currentRouteLine = null;

const greenIcon = new L.Icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-green.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
});

const redIcon = new L.Icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
});

map.on('click', function(e) {
    const lat = e.latlng.lat;
    const lon = e.latlng.lng;

    if (!startMarker) {
        startMarker = L.marker([lat, lon], {icon: greenIcon}).addTo(map).bindPopup("Start").openPopup();
        console.log(`Start selected: ${lat}, ${lon}`);
    } else if (!endMarker) {
        endMarker = L.marker([lat, lon], {icon: redIcon}).addTo(map).bindPopup("End").openPopup();
        console.log(`End selected: ${lat}, ${lon}`);

        const startLat = startMarker.getLatLng().lat;
        const startLon = startMarker.getLatLng().lng;
        const endLat = endMarker.getLatLng().lat;
        const endLon = endMarker.getLatLng().lng;

        const url = `http://localhost:8080/route?start_lat=${startLat}&start_lon=${startLon}&end_lat=${endLat}&end_lon=${endLon}`;
        console.log(`Fetching route from: ${url}`);

        fetch(url)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`Status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                console.log("Route data received successfully!", data);

                if (currentRouteLine) {
                    map.removeLayer(currentRouteLine);
                }

                if (!data.path || data.path.length === 0) {
                    console.warn("Backend returned an empty path array!");
                    return;
                }

                const latLngs = data.path.map(point => [point.lat, point.lon]);

                currentRouteLine = L.polyline(latLngs, {
                    color: '#2A81CB',
                    weight: 5,
                    opacity: 0.85
                }).addTo(map);

                map.fitBounds(currentRouteLine.getBounds(), { padding: [50, 50] });
            })
            .catch(error => {
                console.error("Fetch failed:", error);
            });
    } else {
        map.removeLayer(startMarker);
        map.removeLayer(endMarker);
        if (currentRouteLine) {
            map.removeLayer(currentRouteLine);
        }
        startMarker = null;
        endMarker = null;
        currentRouteLine = null;
        console.log("Markers and routes reset");
    }
});