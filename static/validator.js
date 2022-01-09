const passwordMatch = (p1,p2) => {
    return p1 === p2;
}

const isSlalomEmail = (email)=>{
    const re = /\w+.?\w+\@slalom.com/gm;
    return !!email.match(re);
}