function geoFindMe(map, marker) {

  function success(position) {
    const latitude  = position.coords.latitude;
    const longitude = position.coords.longitude;
    marker.setLatLng([latitude, longitude]);
    map.setView([latitude,longitude], 12);
    const lat = document.getElementById('lat');
    const lon = document.getElementById('lon');
    lat.value = latitude;
    lon.value = longitude;
    console.log(`moving to Latitude: ${latitude} °, Longitude: ${longitude} °`);
  }

  function error() {
    console.warn('Unable to retrieve your location');
  }

  if(!navigator.geolocation) {
    console.warn('Geolocation is not supported by your browser');
  } else {
    console.log('Locating…');
    navigator.geolocation.getCurrentPosition(success, error);
  }

}
