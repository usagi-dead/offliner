export const getWishList = () => {
    return JSON.parse(localStorage.getItem("wishList")) || [];
}

export const addToWishList = (item) => {
    const wishList = getWishList();
    if(!wishList.some(i => i.name === item.name)){
        wishList.push(item);
        localStorage.setItem("wishList", JSON.stringify(wishList));
    }
}

export const removeFromWishList = (itemName) => {
    let wishList = getWishList();
    wishList = wishList.filter(item => item.name !== itemName);
    localStorage.setItem("wishList", JSON.stringify(wishList));
}
