async function loadTimeline() {
    const timelineDiv = document.getElementById('timeline');
    timelineDiv.innerHTML = '<div class="loading">Загрузка...</div>';
    
    try {
        const response = await apiRequest('/timeline');
        if (!response.ok) throw new Error('Failed to load timeline');
        
        const data = await response.json();
        
        if (!data.tweets || data.tweets.length === 0) {
            timelineDiv.innerHTML = '<div class="empty-state">📭 Нет твитов. Подпишитесь на кого-нибудь!</div>';
            return;
        }
        
        timelineDiv.innerHTML = data.tweets.map(tweet => renderTweet(tweet)).join('');
        
        // Привязываем обработчики после рендера
        document.querySelectorAll('.like-btn').forEach(btn => {
            btn.addEventListener('click', () => likeTweet(parseInt(btn.dataset.tweetId)));
        });
        
        document.querySelectorAll('.ai-improve-btn').forEach(btn => {
            btn.addEventListener('click', () => showAlternatives(btn.dataset.tweetText));
        });
        
    } catch (error) {
        timelineDiv.innerHTML = '<div class="empty-state">⚠️ Ошибка загрузки ленты</div>';
        console.error(error);
    }
}

function renderTweet(tweet) {
    const author = tweet.author || {};
    const avatarLetter = author.username ? author.username.charAt(0).toUpperCase() : '?';
    const timeAgo = formatTimeAgo(tweet.created_at);
    
    return `
        <div class="tweet" data-tweet-id="${tweet.id}">
            <div class="tweet-header">
                <div class="tweet-avatar">${avatarLetter}</div>
                <div class="tweet-info">
                    <span class="tweet-author">${author.username || 'Unknown'}</span>
                    <span class="tweet-time">${timeAgo}</span>
                </div>
            </div>
            <div class="tweet-text">${escapeHtml(tweet.text)}</div>
            <div class="tweet-actions">
                <div class="tweet-action like-btn" data-tweet-id="${tweet.id}">
                    ❤️ ${tweet.like_count || 0}
                </div>
                <div class="tweet-action">
                    💬 ${tweet.reply_count || 0}
                </div>
            </div>
        </div>
    `;
}

function formatTimeAgo(timestamp) {
    const date = new Date(timestamp);
    const now = new Date();
    const diff = Math.floor((now - date) / 1000);
    
    if (diff < 60) return 'только что';
    if (diff < 3600) return `${Math.floor(diff / 60)} мин назад`;
    if (diff < 86400) return `${Math.floor(diff / 3600)} ч назад`;
    return date.toLocaleDateString();
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

async function likeTweet(tweetId) {
    // TODO: Implement like functionality
    console.log('Like tweet:', tweetId);
}

// Загрузка твитов пользователя по имени
async function getTweetsByUsername(username, container) {
    if (!container) return;
    
    container.innerHTML = '<div class="loading">Загрузка твитов...</div>';
    
    try {
        const response = await apiRequest(`/users/${username}/tweets`);
        if (!response.ok) throw new Error('Failed to load tweets');
        
        const data = await response.json();
        const tweets = data.tweets || data;
        
        if (!tweets || tweets.length === 0) {
            container.innerHTML = '<div class="empty-state">📭 У пользователя нет твитов</div>';
            return;
        }
        
        container.innerHTML = tweets.map(tweet => renderTweet(tweet)).join('');
        
        // Привязываем обработчики
        document.querySelectorAll('.like-btn').forEach(btn => {
            btn.addEventListener('click', () => likeTweet(parseInt(btn.dataset.tweetId)));
        });
        
    } catch (error) {
        console.error('Error loading user tweets:', error);
        container.innerHTML = '<div class="empty-state">Ошибка загрузки твитов</div>';
    }
}