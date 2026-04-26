function showError(msg) {
  var el = document.getElementById("auth-error");
  el.textContent = msg;
  el.style.display = "block";
}

function hideError() {
  var el = document.getElementById("auth-error");
  el.style.display = "none";
}

function setLoading(btn, loading) {
  btn.disabled = loading;
  btn.textContent = loading ? "Please wait..." : btn.dataset.label;
}

function handleSignUp(e) {
  e.preventDefault();
  hideError();

  var username = document.getElementById("username").value.trim();
  var password = document.getElementById("password").value;
  var repeatedPassword = document.getElementById("repeated-password").value;

  if (username === "" || password === "") {
    showError("Username and password are required.");
    return;
  }
  if (password !== repeatedPassword) {
    showError("Passwords do not match.");
    return;
  }

  var btn = e.target.querySelector("button[type=submit]");
  setLoading(btn, true);

  fetch("/signup", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      Username: username,
      Password: password,
      RepeatedPassword: repeatedPassword,
    }),
  })
    .then(function (res) {
      if (!res.ok) {
        return res.text().then(function (text) {
          throw new Error(text || "Sign up failed");
        });
      }
      return res.json();
    })
    .then(function (data) {
      localStorage.setItem("token", data.token);
      window.location.href = "/";
    })
    .catch(function (err) {
      showError(err.message);
      setLoading(btn, false);
    });
}

function handleSignIn(e) {
  e.preventDefault();
  hideError();

  var username = document.getElementById("username").value.trim();
  var password = document.getElementById("password").value;

  if (username === "" || password === "") {
    showError("Username and password are required.");
    return;
  }

  var btn = e.target.querySelector("button[type=submit]");
  setLoading(btn, true);

  fetch("/signin", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      Username: username,
      Password: password,
    }),
  })
    .then(function (res) {
      if (!res.ok) {
        return res.text().then(function (text) {
          throw new Error(text || "Sign in failed");
        });
      }
      return res.json();
    })
    .then(function (data) {
      localStorage.setItem("token", data.token);
      window.location.href = "/";
    })
    .catch(function (err) {
      showError(err.message);
      setLoading(btn, false);
    });
}
