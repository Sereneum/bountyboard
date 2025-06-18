// modal.js - Управление модальными окнами
function setupModals() {
    const modals = document.querySelectorAll('.modal');
    const openButtons = document.querySelectorAll('.btn-desc');
    const closeButtons = document.querySelectorAll('.close');

    openButtons.forEach((btn, idx) => {
        btn.addEventListener('click', () => modals[idx].classList.add('show'));
    });

    closeButtons.forEach((btn, idx) => {
        btn.addEventListener('click', () => modals[idx].classList.remove('show'));
    });

    window.addEventListener('click', (e) => {
        modals.forEach(modal => {
            if (e.target === modal) modal.classList.remove('show');
        });
    });
}

// task.js - Работа с задачами
async function completeTask(btn) {
    const taskItem = btn.closest(".task-item");
    const taskId = taskItem.dataset.id || "unknown";

    try {
        await fetch("/complete-task", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ id: taskId }),
        });

        taskItem.classList.add("done");
        for (let i = 0; i < 16; i++) createCoin(taskItem);
    } catch (err) {
        console.error("Ошибка отправки:", err);
    }
}

// effects.js - Визуальные эффекты
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

// Инициализация
document.addEventListener("DOMContentLoaded", () => {
    setupModals();

    document.querySelectorAll(".btn-complete").forEach(btn => {
        btn.addEventListener("click", () => completeTask(btn));
    });
});