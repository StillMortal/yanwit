// Функция публикации твита
async function postTweet() {
    const text = document.getElementById('tweet-text').value;
    if (!text.trim()) return;
    
    try {
        const response = await apiRequest('/tweets', {
            method: 'POST',
            body: JSON.stringify({ text })
        });
        
        if (!response.ok) throw new Error('Failed to post tweet');
        
        document.getElementById('tweet-text').value = '';
        document.getElementById('char-counter').textContent = '280';
        loadTimeline();
        
    } catch (error) {
        alert('Ошибка публикации: ' + error.message);
    }
}

// Функция обновления счётчика символов
function updateCharCounter() {
    const textarea = document.getElementById('tweet-text');
    const counter = document.getElementById('char-counter');
    if (!textarea || !counter) return;
    
    const remaining = 280 - textarea.value.length;
    counter.textContent = remaining;
    counter.classList.remove('warning', 'danger');
    
    if (remaining < 20) counter.classList.add('warning');
    if (remaining < 0) counter.classList.add('danger');
}

// Функция проверки манипуляции (если не определена в ai.js)
if (typeof window.checkManipulation === 'undefined') {
    window.checkManipulation = async function() {
        const text = document.getElementById('tweet-text').value;
        if (!text.trim()) return;
        
        const warningDiv = document.getElementById('manipulation-warning');
        if (!warningDiv) return;
        
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
                warningDiv.style.display = 'block';
            } else {
                warningDiv.style.display = 'none';
            }
            
            setTimeout(() => {
                warningDiv.style.display = 'none';
            }, 5000);
            
        } catch (error) {
            console.error('Error checking manipulation:', error);
        }
    };
}

// Привязка кнопок после загрузки DOM
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, binding buttons...');
    
    // Кнопка Проверить манипуляцию
    const checkBtn = document.getElementById('check-manipulation-btn');
    if (checkBtn) {
        // Удаляем старые обработчики
        checkBtn.removeAttribute('onclick');
        // Добавляем новый
        checkBtn.addEventListener('click', function(e) {
            e.preventDefault();
            window.checkManipulation();
        });
        console.log('Check button bound');
    } else {
        console.error('Check button not found!');
    }
    
    // Кнопка AI улучшить
    const improveBtn = document.getElementById('improve-tweet-btn');
    if (improveBtn && typeof window.improveDraft === 'function') {
        improveBtn.addEventListener('click', function(e) {
            e.preventDefault();
            window.improveDraft();
        });
        console.log('Improve button bound');
    }
    
    // Кнопка публикации
    const postBtn = document.getElementById('post-tweet-btn');
    if (postBtn) {
        postBtn.addEventListener('click', function(e) {
            e.preventDefault();
            postTweet();
        });
        console.log('Post button bound');
    }
    
    // Счётчик символов
    const textarea = document.getElementById('tweet-text');
    if (textarea) {
        textarea.addEventListener('input', updateCharCounter);
    }
});