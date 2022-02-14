
const createUser = (email, user, password, success, fail) => {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    const payload = JSON.stringify({
        "email": email,
        "user": user,
        "pass": password
    });

    const requestOptions = {
        method: 'POST',
        headers: myHeaders,
        body: payload,
        redirect: 'follow'
    };

    fetch("/user", requestOptions)
        .then(response =>{
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) })
            }
            return response.text()
        })
        .then(result => success(result))
        .catch(error => fail(error));
}

const fetchAllSlalomers = (cookieval, success, fail)=>{
    const myHeaders = new Headers();

    myHeaders.append("Cookie", "mysession=" + cookieval);
    
    const requestOptions = {
        method: 'GET',
        headers: myHeaders,
        redirect: 'follow'
    };
    
    fetch("/users", requestOptions)
        .then(response =>{
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) })
            }
            return response.text()
        })
        .then(result => success(result))
        .catch(error => fail(error));
}

const updateUser = (cookie, user, success,fail) =>{
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Cookie", "mysession=" + cookie);

    var requestOptions = {
        method: 'PUT',
        headers: myHeaders,
        body: JSON.stringify(user),
        redirect: 'follow'
    };

    fetch("/user", requestOptions)
        .then(response =>{
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) })
            }
            return response.text()
        })
        .then(result => success(result))
        .catch(error => fail(error));
}

const LoginUser = (email, user, pass, success,fail) =>{
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    var raw = JSON.stringify({
        "email": email,
        "user": user,
        "pass": pass
      });

    var requestOptions = {
        method: 'PUT',
        headers: myHeaders,
        body: raw,
        redirect: 'follow'
    };

    fetch("/login", requestOptions)
        .then(response =>{
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) })
            }
            return response.text()
        })
        .then(result => success(result))
        .catch(error => fail(error));
}

const LogOut = (cookie, success,fail) =>{
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Cookie", cookie);


    var requestOptions = {
        method: 'DELETE',
        headers: myHeaders,
        body: {},
        redirect: 'follow'
    };

    fetch("/login", requestOptions)
        .then(response =>{
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) })
            }
            return response.text()
        })
        .then(result => success(result))
        .catch(error => fail(error));
}