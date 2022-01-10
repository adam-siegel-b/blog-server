
const createUser = (email, user, password, success, fail) => {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    const payload = JSON.stringify({
        "email": email,
        "user": user,
        "password": password
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