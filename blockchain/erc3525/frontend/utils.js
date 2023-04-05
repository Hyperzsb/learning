function toggleModal(title, raw) {
  document.getElementById("uri-modal-label").innerText = title;
  // Display the SVG image
  // Convert the string to a Blob object
  const svgBlob = new Blob([raw], {
    type: "image/svg+xml;charset=utf-8",
  });
  // Create a URL for the Blob object
  const svgUrl = URL.createObjectURL(svgBlob);
  // Create an <img> element and set its src attribute to the URL
  const img = document.createElement("img");
  img.src = svgUrl;
  img.style.width = "100%";
  img.style.border = "1px solid black";

  // Append the <img> element to a modal to display it
  const container = document.getElementById("uri-container");
  while (container.firstChild) {
    container.removeChild(container.firstChild);
  }
  container.appendChild(img);

  // Activate the modal
  const modal = new bootstrap.Modal("#uri-modal", {});
  modal.show();
}

function toggleToast(status, body) {
  const toast = new bootstrap.Toast("#toast", {});

  if (status == true) {
    document.getElementById("toast-title-failure").style.display = "none";
    document.getElementById("toast-title-success").style.display =
      "inline-block";
  } else {
    document.getElementById("toast-title-failure").style.display =
      "inline-block";
    document.getElementById("toast-title-success").style.display = "none";
  }

  const now = new Date();
  const hours = now.getHours();
  const minutes = now.getMinutes();
  const seconds = now.getSeconds();
  const timeString = `${hours}:${minutes}:${seconds}`;
  document.getElementById("toast-time").innerText = timeString;

  document.getElementById("toast-body").innerText = body;
  toast.show();
}