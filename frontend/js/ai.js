// Глобальные функции для AI кнопок

window.improveDraft = async function() {
    const text = document.getElementById('tweet-text').value;
    if (!text.trim()) {
        alert('Напишите текст твита перед улучшением');
        return;
    }
    
    const modal = document.getElementById('alternatives-modal');
    const listDiv = document.getElementById('alternatives-list');
    
    listDiv.innerHTML = '<div class="loading">Генерация вариантов...</div>';
    modal.style.display = 'flex';
    
    try {
        const response = await apiRequest('/ai/alternatives', {
            method: 'POST',
            body: JSON.stringify({ text: text, count: 3, style: 'funny' })
        });
        
        const data = await response.json();
        
        // Нормализуем каждый вариант перед отображением
        const normalizedAlternatives = data.alternatives.map(alt => normalizeText(alt));
        
        listDiv.innerHTML = normalizedAlternatives.map(alt => `
            <div class="alternative-item" onclick="window.replaceDraftText('${escapeHtml(alt)}')">
                ${escapeHtml(alt)}
            </div>
        `).join('');
        
    } catch (error) {
        listDiv.innerHTML = '<div class="empty-state">Ошибка генерации вариантов</div>';
        console.error(error);
    }
};

// Функция нормализации текста (исправляет заглавные буквы после запятой)
function normalizeText(text) {
    if (!text) return text;
    const parts = text.split(/(,)/);
    for (let i = 1; i < parts.length; i++) {
        if (parts[i] === ',' && i + 1 < parts.length) {
            const nextPart = parts[i + 1];
            if (nextPart.length > 0 && nextPart[0] === ' ') {
                parts[i + 1] = nextPart.toLowerCase();
            }
        }
    }
    return parts.join('');
}

window.replaceDraftText = function(newText) {
    document.getElementById('tweet-text').value = newText;
    document.getElementById('alternatives-modal').style.display = 'none';
    if (typeof updateCharCounter === 'function') {
        updateCharCounter();
    }
};

window.checkManipulation = async function() {
    console.log('checkManipulation вызвана');
    
    const text = document.getElementById('tweet-text').value;
    if (!text.trim()) {
        console.log('Текст пуст');
        return;
    }
    
    const warningDiv = document.getElementById('manipulation-warning');
    if (!warningDiv) {
        console.error('warningDiv не найден');
        return;
    }
    
    // Принудительно показываем
    warningDiv.style.setProperty('display', 'block', 'important');
    warningDiv.innerHTML = '<div class="warning-title">⏳ Проверка текста...</div>';
    
    try {
        const response = await apiRequest('/ai/detect-manipulation', {
            method: 'POST',
            body: JSON.stringify({ text })
        });
        
        const data = await response.json();
        
        if (data.has_manipulation) {
            warningDiv.innerHTML = `
                <div class="warning-title">⚠️ Обнаружена манипуляция: ${data.types.join(', ')}</div>
                <div>Уверенность: ${Math.round(data.confidence * 100)}%</div>
                <div class="warning-suggestion">💡 ${data.suggestions.join(' • ')}</div>
            `;
        } else {
            warningDiv.innerHTML = '<div class="warning-title">✅ Манипуляций не обнаружено</div>';
        }
        warningDiv.style.setProperty('display', 'block', 'important');
        
        // Скрываем сообщение через 5 секунд
        setTimeout(() => {
            warningDiv.style.setProperty('display', 'none', 'important');
        }, 5000);
        
    } catch (error) {
        console.error('Ошибка:', error);
        warningDiv.innerHTML = `<div class="warning-title">❌ Ошибка: ${error.message}</div>`;
        warningDiv.style.setProperty('display', 'block', 'important');
        
        // Скрываем сообщение об ошибке через 5 секунд
        setTimeout(() => {
            warningDiv.style.setProperty('display', 'none', 'important');
        }, 5000);
    }
};

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Инициализация модального окна
document.addEventListener('DOMContentLoaded', () => {
    const modal = document.getElementById('alternatives-modal');
    const closeBtn = document.querySelector('.modal-close');
    
    if (closeBtn) {
        closeBtn.addEventListener('click', () => {
            if (modal) modal.style.display = 'none';
        });
    }
    
    if (modal) {
        window.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.style.display = 'none';
            }
        });
    }
});