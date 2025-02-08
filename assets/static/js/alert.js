export const showAlert = (msj) => {
    const alertBox = document.getElementById("customAlert");
    alertBox.style.display = "block";
    alertBox.textContent=msj

    setTimeout(() => {
      alertBox.style.display = "none";
    }, 5000);
}