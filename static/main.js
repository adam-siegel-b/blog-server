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

/**
 * 
 * @param {number} lat in degrees as number
 * @param {number} lon in degrees as number
 * @param {object} map object for current context
 * @returns object
 */
const renderCurrentUserMarker = (lat, lon, map) => {
    let mymarker = new L.Marker([lat,lon] , iconOptions);
    mymarker.bindTooltip("move me to where you want to be").openTooltip();
    mymarker.addTo(map);
    // center on my current location
    geoFindMe(map, mymarker);
    // My marker listener
    mymarker.on('dragend', (event) => {
        let latlng = event.target.getLatLng();
        updateLatLon(latlng.lat, latlng.lng);
        });
    return mymarker;
}

/**
 * 
 * @param {number} lat 
 * @param {number} lon 
 * @param {object} map 
 * @param {string} desc 
 * @param {string} icon 
 * @param {string} user 
 * @returns object
 */
const renderUserMarker = (lat, lon, map, desc, icon, user) => {
    const iconM = {
        iconUrl:'https://i.ibb.co/k4W53P7/location.png',
        iconSize:[37,37]
       }
    if (icon.length){
        iconM.iconUrl = icon
    }
    let uicon = L.icon(iconM);
    let Options = {
        title:user,
        alt: desc,
        opacity: 0.75,
        draggable:false,
        icon:uicon
    }
    let mymarker = new L.Marker([lat,lon], Options);
    mymarker.bindTooltip(desc);
    mymarker.addTo(map);
    return mymarker;
}

const getOtherSlalomers = ()=>{
    const success = (res)=>{
        const response = JSON.parse(res);
        console.log(response);
        response.Users.map(val => renderUserMarker(val.location.lat, val.location.lon, window.map, val.user, val.photo, val.id));
    }
    const fail = (err)=>{
        console.warn(err);
    }
    fetchAllSlalomers(getCookie('mysession'),success, fail);
}

/**
 * 
 * @param {boolean} isLoggedIn 
 * @returns void
 */
const setLoggedIn = (isLoggedIn)=>{
    let signInToggle = document.getElementById('loginup');
    let editToggle = document.getElementById('edituser');
    if(isLoggedIn){
        document.querySelector('.form.modal').classList.add("hide");
        signInToggle.classList.add('hide');
        editToggle.classList.remove('hide');
        if (window.map){
            getOtherSlalomers();
        }
        return;
    }
    document.querySelector('.form.modal').classList.remove("hide");
    signInToggle.classList.remove('hide');
    editToggle.classList.add('hide');
}

/**
 * 
 * @param {string} id 
 * @param {string} value 
 * @returns 
 */
const updateField = (id, value)=>{
    const field = document.getElementById(id);
    if (typeof field.value === 'string'){
        field.value = value;
        return false;
    }
    return true;
}

const updateLogin= (res)=>{
    const response = JSON.parse(res);
    console.log(response);
    if (response.location && response.user && response.email){
        let marker = renderCurrentUserMarker(response.location.lat,response.location.lon, window.map);
        updateField('lat',response.location.lat);
        updateField('lon',response.location.lon);
        updateField('placename', response.location['loc-name']);
        updateField('displayname', response.user);
        updateField('accountemail', response.email);
        updateField('userid',response.id);
        if (!response.location.lat && !response.location.lon){
            geoFindMe(window.map, marker);
        }
    }
    setLoggedIn(checkLogin());
}

const main = () => {
    // Initial draw
    let map = L.map('map').setView(LA, 12);
    // check if the user is logged in
    setLoggedIn(checkLogin());
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

    // signup & login 
    let loginToggle = document.getElementById('signintoggle');
    loginToggle.addEventListener('click', (event)=>{
        event.preventDefault();
        document.querySelector('.signin').classList.remove("hide");
        document.querySelector('.signup').classList.add("hide");
        document.getElementById('signuptoggle').classList.remove("hide")
        event.target.classList.add("hide");
    });

    let signupToggle = document.getElementById('signuptoggle');
    signupToggle.addEventListener('click', (event)=>{
        event.preventDefault();
        document.querySelector('.signin').classList.add("hide");
        document.querySelector('.signup').classList.remove("hide");
        document.getElementById('signintoggle').classList.remove("hide")
        event.target.classList.add("hide");
    });

    // login handler
    const login = document.getElementById('login');
    login.addEventListener('click', (event)=>{
        const user = document.getElementById('login-username');
        const password = document.getElementById('login-password');
        const errorMsg = document.getElementById('login-errors');

        if (!user.value.length){
            user.classList.add("error");
            errorMsg.textContent += "User can't be blank. ";
        }else{
            user.classList.remove("error");
            errorMsg.textContent.replace("Bad username or password. ","");
            errorMsg.textContent.replace("User can't be blank. ","");
        }
        if (!password.value.length){
            password.classList.add("error");
            errorMsg.textContent += "password can't be blank. ";
        }else{
            password.classList.remove("error");
            errorMsg.textContent.replace("Bad username or password. ","");
            errorMsg.textContent.replace("password can't be blank. ","");
        }

        const fields = [user,password];
        let errorCount = 0;
        fields.map((item)=>{if (item.classList.contains("error")){errorCount++}});
        if (errorCount){
            return;
        }

        const handleError = (err) => {
            const error = JSON.parse(err.message);
            console.warn(error);
            user.classList.add("error");
            password.classList.add("error");
            errorMsg.textContent += "Bad username or password. ";
        }
        let usr = user.value;
        let email = ''; 
        if (user.value.indexOf('@slalom.com') > -1){
            email = user.value;
            usr = '';
        }
        LoginUser(email, usr, password.value, updateLogin, handleError);
    });

    // signup handler
    const signup = document.getElementById('signup');
    signup.addEventListener('click', (event)=>{
        event.preventDefault();
        const user = document.getElementById('signup-user');
        const email = document.getElementById('signup-email');
        const password = document.getElementById('signup-password');
        const password2 = document.getElementById('signup-password2');
        const errorMsg = document.getElementById('signup-errors');
        if (!user.value.length){
            user.classList.add("error");
            errorMsg.textContent += "User can't be blank. ";
        }else{
            user.classList.remove("error");
            errorMsg.textContent.replace("User can't be blank. ","");
            errorMsg.textContent.replace("Username already taken. ","");
        }
        if (!isSlalomEmail(email.value)){
            email.classList.add("error");
            errorMsg.textContent += "Not a valid Slalom email. ";
        }else{
            email.classList.remove("error");
            errorMsg.textContent.replace("Not a valid Slalom email. ","");
            errorMsg.textContent.replace("Email used already. ","");
        }
        if (!passwordMatch(password.value, password2.value)){
            password2.classList.add("error");
            password.classList.add("error");
            errorMsg.textContent += "Passwords must match. ";
        }else{
            password2.classList.remove("error");
            password.classList.remove("error");
            errorMsg.textContent.replace("Not a valid Slalom email. ","");
        }
        const fields = [user,email,password2,password];
        let errorCount = 0;
        fields.map((item)=>{if (item.classList.contains("error")){errorCount++}});
        if (errorCount){
            return;
        }
        const handleError = (err) => {
            const error = JSON.parse(err.message);
            console.warn(error);
            if(error.value.indexOf("slalomer_name_key") > -1){
                user.classList.add("error");
                errorMsg.textContent += "Username already taken. ";
            }
            if(error.value.indexOf("slalomer_email_key") > -1){
                email.classList.add("error");
                errorMsg.textContent += "Email used already. ";
            }
        }
        createUser(email.value,user.value,user.password, updateLogin, handleError);
    });

    let upUser = document.getElementById('updateUser');
    upUser.addEventListener('click',(event)=>{
        event.preventDefault();
        const displayname = document.getElementById('displayname');
        const accountemail = document.getElementById('accountemail');
        const placename = document.getElementById('placename');
        const lat = document.getElementById('lat');
        const lon = document.getElementById('lon');
        const userid = document.getElementById('userid');

        const errorMsg = document.getElementById('edit-error');
        const payload = {
            "id": userid.value,
            "email": accountemail.value,
            "user": displayname.value,
            "location":{
                "lat": parseFloat(lat.value),
                "lon": parseFloat(lon.value),
                "loc-name": placename.value
            }
        };
        const success = (val) =>{
            const result = JSON.parse(val);
            console.log(result);
            displayname.value = result.user;
            accountemail.value = result.email;
            placename.value = result.location['loc-name'];
            lat.value = result.location.lat;
            lon.value = result.location.lon;
            userid.value = result.id;
        }
        // TODO: error on duplicate email/username
        const fail = (error) => {
            const errorMsg = document.getElementById('edit-error');
            console.warn(error);
            errorMsg.textContent = error;
        }
        // TODO: validate lat/lon
        const fields = [displayname, accountemail,placename,lat,lon, userid];
        let errorCount = 0;
        fields.map((item)=>{if (item.classList.contains("error")){errorCount++}});
        if (errorCount){
            return;
        }
        updateUser(getCookie('mysession'),payload, success, fail);
    });

    // main menu toggles
    let signInToggle = document.getElementById('loginup');
    signInToggle.addEventListener('click', (event)=>{
        event.preventDefault();
        document.querySelector('.user.modal').classList.add("hide");
        document.querySelector('.form.modal').classList.remove("hide");
    });

    let editToggle = document.getElementById('edituser');
    editToggle.addEventListener('click', (event)=>{
        event.preventDefault();
        document.querySelector('.form.modal').classList.add("hide");
        document.querySelector('.user.modal').classList.remove("hide");
    });
    
    // logout handler
    const signout = document.getElementById('signout');
    signout.addEventListener('click', (event)=>{
        event.preventDefault();
        LogOut(document.cookie,()=> setLoggedIn(checkLogin()),() => setLoggedIn(checkLogin()));
    });

    // Set up the OSM layer
    L.tileLayer(
        'https://cartodb-basemaps-{s}.global.ssl.fastly.net/dark_all/{z}/{x}/{y}.png'
    ).addTo(map);

    // attach map to the window object for maximum hackyness
    window.map = map;
}
main();