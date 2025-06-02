function fetchUsers() {
  fetch("http://localhost:5500/users")
    .then((response) => {
      if (!response.ok) throw new Error("Network response was not ok");
      return response.json();
    })
    .then((users) => {
      const list = document.getElementById("user-list");
      list.innerHTML = "";
      users.forEach((user) => {
        const button = document.createElement("button");
        button.textContent = user.name;
        button.onclick = () => fetchUserById(user.id);
        list.appendChild(button);
      });
    })
    .catch((error) => {
      alert("Error: " + error.message);
    });
}

function fetchUserById(id) {
  fetch(`http://localhost:5500/users/${id}`)
    .then((response) => {
      if (!response.ok) throw new Error("User not found");
      return response.json();
    })
    .then((user) => {
      alert(
        `Name: ${user.name}\nAge: ${user.age}\nEmail: ${user.email}\nNationality: ${user.nationality}`
      );
    })
    .catch((error) => {
      alert("Error: " + error.message);
    });
}
