// Document ready
document.addEventListener('DOMContentLoaded', function() {
    // Initialize tooltips
    const tooltipElements = document.querySelectorAll('[data-tooltip]');
    tooltipElements.forEach(el => {
        el.addEventListener('mouseenter', showTooltip);
        el.addEventListener('mouseleave', hideTooltip);
    });
    
    // Add animation class to elements
    const animatedElements = document.querySelectorAll('.animate-on-scroll');
    animatedElements.forEach(el => {
        el.classList.add('opacity-0', 'translate-y-6');
    });
    
    // Check for tasks due soon (every minute)
    checkDueTasks();
    setInterval(checkDueTasks, 60000);
});

function showTooltip(e) {
    const tooltipText = e.target.getAttribute('data-tooltip');
    const tooltip = document.createElement('div');
    tooltip.className = 'absolute z-50 bg-gray-800 text-white text-xs rounded py-1 px-2 whitespace-nowrap';
    tooltip.textContent = tooltipText;
    
    const rect = e.target.getBoundingClientRect();
    tooltip.style.top = `${rect.top - 30}px`;
    tooltip.style.left = `${rect.left + rect.width / 2}px`;
    tooltip.style.transform = 'translateX(-50%)';
    
    tooltip.id = 'current-tooltip';
    document.body.appendChild(tooltip);
}

function hideTooltip() {
    const tooltip = document.getElementById('current-tooltip');
    if (tooltip) {
        tooltip.remove();
    }
}

function checkDueTasks() {
    const now = new Date();
    const taskCards = document.querySelectorAll('.task-card');
    
    taskCards.forEach(card => {
        const dueTimeStr = card.querySelector('.due-time').textContent;
        const dueTime = new Date(dueTimeStr);
        
        // If task is due within 30 minutes and not completed
        if (!card.classList.contains('completed') && dueTime - now < 30 * 60 * 1000) {
            card.classList.add('border-red-500', 'animate-pulse');
            
            // Show notification if not already shown
            if (!card.hasAttribute('data-notified')) {
                showNotification(`Task "${card.querySelector('h3').textContent}" is due soon!`);
                card.setAttribute('data-notified', 'true');
            }
        }
    });
}

function showNotification(message) {
    // Check if browser supports notifications
    if (!("Notification" in window)) {
        console.log("This browser does not support desktop notification");
        return;
    }
    
    // Check if notification permissions have already been granted
    if (Notification.permission === "granted") {
        new Notification(message);
    } else if (Notification.permission !== "denied") {
        // Otherwise, ask the user for permission
        Notification.requestPermission().then(permission => {
            if (permission === "granted") {
                new Notification(message);
            }
        });
    }
}