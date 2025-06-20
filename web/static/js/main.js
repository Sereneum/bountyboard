document.addEventListener("DOMContentLoaded", () => {
    setupModals();

    document.querySelectorAll(".btn-complete").forEach(btn => {
        btn.addEventListener("click", () => completeTask(btn));
    });
});

function setupModals() {
    const openButtons = document.querySelectorAll('.btn-desc');
    const closeButtons = document.querySelectorAll('.close');

    openButtons.forEach((btn) => {
        const modalId = btn.dataset.target;
        btn.addEventListener('click', () => {
            document.getElementById(modalId).classList.add('show');
        });
    });

    closeButtons.forEach((btn) => {
        const modalId = btn.dataset.target;
        btn.addEventListener('click', () => {
            document.getElementById(modalId).classList.remove('show');
        });
    });

    window.addEventListener('click', (e) => {
        if (e.target.classList.contains('modal')) {
            e.target.classList.remove('show');
        }
    });
}

async function completeTask(btn) {
    const taskItem = btn.closest(".task-item");
    const taskId = btn.dataset.id;

    try {

        await fetch("/complete-task", { // 🔗 Backend: обработай POST /complete-task
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ id: taskId }),
        });

        taskItem.classList.add("done");
        for (let i = 0; i < 12; i++) createCoin(taskItem);
    } catch (err) {
        console.error("Failed to complete task:", err);
    }
}

function createCoin(parent) {
    const coin = document.createElement("div");
    coin.classList.add("coin-burst");

    const rect = parent.getBoundingClientRect();
    const x = Math.random() * rect.width;
    const y = Math.random() * rect.height;

    coin.style.left = `${x}px`;
    coin.style.top = `${y}px`;

    parent.appendChild(coin);
    setTimeout(() => coin.remove(), 1200);
}

document.querySelectorAll(".btn-delete").forEach(btn => {
    btn.addEventListener("click", () => deleteTask(btn));
});

async function deleteTask(btn) {
    const taskItem = btn.closest(".task-item");
    const taskId = btn.dataset.id;

    try {
        await fetch("/delete-task", { // 🔗 Backend: обработай POST /delete-task
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ id: taskId }),
        });

        taskItem.remove();
    } catch (err) {
        console.error("Failed to delete task:", err);
    }
}
