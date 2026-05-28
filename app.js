const map = L.map('map').setView([39.15, -75.52], 9);

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '© OpenStreetMap contributors'
}).addTo(map);

let startMarker = null;
let endMarker = null;

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
    } else {
        map.removeLayer(startMarker);
        map.removeLayer(endMarker);
        startMarker = null;
        endMarker = null;
        console.log("Markers reset");
    }
});