const taskList = document.getElementById("task-list");
const taskInput = document.getElementById("task-input");
const addBtn = document.getElementById("add-btn");

function createTaskElement(todo) {
  const li = document.createElement("li");
  li.className = "task-item";
  li.dataset.id = todo.id;

  const span = document.createElement("span");
  span.textContent = todo.text;

  span.addEventListener("dblclick", function () {
    startEditing(li, todo);
  });

  const editBtn = document.createElement("button");
  editBtn.className = "edit-btn";
  editBtn.textContent = "\u270E";
  editBtn.addEventListener("click", function () {
    startEditing(li, todo);
  });

  const deleteBtn = document.createElement("button");
  deleteBtn.className = "delete-btn";
  deleteBtn.textContent = "\u2715";
  deleteBtn.addEventListener("click", function () {
    deleteTodo(todo.id, li);
  });

  const actions = document.createElement("div");
  actions.className = "task-actions";
  actions.appendChild(editBtn);
  actions.appendChild(deleteBtn);

  li.appendChild(span);
  li.appendChild(actions);
  return li;
}

function loadTodos() {
  fetch("/all/")
    .then(function (res) {
      if (!res.ok) throw new Error("failed to fetch");
      return res.json();
    })
    .then(function (todos) {
      taskList.innerHTML = "";
      if (!todos || todos.length === 0) return;
      todos.forEach(function (todo) {
        taskList.appendChild(createTaskElement(todo));
      });
    })
    .catch(function (err) {
      console.error("load error:", err);
    });
}

function addTodo() {
  const text = taskInput.value.trim();
  if (text === "") return;

  fetch("/create", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ text: text }),
  })
    .then(function (res) {
      if (!res.ok) throw new Error("failed to create");
      return res.json();
    })
    .then(function (todo) {
      taskList.appendChild(createTaskElement(todo));
      taskInput.value = "";
      taskInput.focus();
    })
    .catch(function (err) {
      console.error("create error:", err);
    });
}

function startEditing(li, todo) {
  if (li.querySelector(".edit-input")) return;

  const span = li.querySelector("span");
  const actions = li.querySelector(".task-actions");
  span.style.display = "none";
  actions.style.display = "none";

  const input = document.createElement("input");
  input.type = "text";
  input.className = "edit-input";
  input.value = todo.text;
  li.insertBefore(input, span);
  input.focus();

  function finishEdit() {
    var newText = input.value.trim();
    if (newText === "" || newText === todo.text) {
      input.remove();
      span.style.display = "";
      actions.style.display = "";
      return;
    }
    updateTodo(todo.id, newText, li, span, input, actions, todo);
  }

  input.addEventListener("keydown", function (e) {
    if (e.key === "Enter") finishEdit();
    if (e.key === "Escape") {
      input.remove();
      span.style.display = "";
      actions.style.display = "";
    }
  });

  input.addEventListener("blur", finishEdit);
}

function updateTodo(id, newText, li, span, input, actions, todo) {
  fetch("/update", {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ id: id, text: newText }),
  })
    .then(function (res) {
      if (!res.ok) throw new Error("failed to update");
      return res.json();
    })
    .then(function (updated) {
      todo.text = updated.text;
      span.textContent = updated.text;
      input.remove();
      span.style.display = "";
      actions.style.display = "";
    })
    .catch(function (err) {
      console.error("update error:", err);
      input.remove();
      span.style.display = "";
      actions.style.display = "";
    });
}

function deleteTodo(id, li) {
  fetch("/delete", {
    method: "DELETE",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ id: id }),
  })
    .then(function (res) {
      if (!res.ok) throw new Error("failed to delete");
      li.remove();
    })
    .catch(function (err) {
      console.error("delete error:", err);
    });
}

addBtn.addEventListener("click", addTodo);

taskInput.addEventListener("keydown", function (e) {
  if (e.key === "Enter") addTodo();
});

loadTodos();
