// set a default for the map
const LA = [34.052235, -118.243683];

//default marker Icon
let customIcon = {
 iconUrl:'https://i.ibb.co/k4W53P7/location.png',
 iconSize:[37,37]
}

let myIcon = L.icon(customIcon);

let iconOptions = {
 title:'me',
 draggable:true,
 icon:myIcon
}

/**
 * 
 * @param {number} lat in degrees as a number 
 * @param {number} lon in degrees as a number
 */
const updateLatLon = (lt,ln) => {
    const lat = document.getElementById('lat');
    const lon = document.getElementById('lon');
    lat.value = lt;
    lon.value = ln;
}


const main = () => {
    // Initial draw
    let map = L.map('map').setView(LA, 12);
    let mymarker = new L.Marker(LA , iconOptions);
    mymarker.bindTooltip("move me to where you want to be").openTooltip()
    mymarker.addTo(map);

    // Find where I'm at
    geoFindMe(map, mymarker);

    // close buttons
    let closeUser = document.getElementById('close-user');
    closeUser.addEventListener('click', (event)=>{
        event.preventDefault();
        document.querySelector('.user.modal').classList.add("hide");
    });

    let closeForm = document.getElementById('close-form');
    closeForm.addEventListener('click', (event)=>{
        event.preventDefault();
        document.querySelector('.form.modal').classList.add("hide");
    });

    // My marker listener
    mymarker.on('dragend', (event) => {
        let latlng = event.target.getLatLng();
        updateLatLon(latlng.lat, latlng.lng);
      });

    // Set up the OSM layer
    L.tileLayer(
        'https://cartodb-basemaps-{s}.global.ssl.fastly.net/dark_all/{z}/{x}/{y}.png'
    ).addTo(map);
}
main();