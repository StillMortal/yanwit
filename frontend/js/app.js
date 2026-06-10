// Главный файл приложения
document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM loaded');
    
    // Проверяем авторизацию
    if (checkAuth && checkAuth() && currentUser) {
        showMainApp();
        if (typeof loadTimeline === 'function') loadTimeline();
    } else {
        showAuthPage();
    }
    
    // Настройка навигации
    setupNavigation();
    
    // Настройка вкладок авторизации
    setupAuthTabs();
    
    // Настройка форм
    setupAuthForms();
    
    // Кнопка выхода
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', () => {
            logout();
        });
    }
    
    // Тёмная тема / Светлая тема
    const themeToggle = document.getElementById('dark-theme-toggle');
    if (themeToggle) {
        // Загружаем сохранённую тему
        const savedTheme = localStorage.getItem('yanwit_theme');
        if (savedTheme === 'light') {
            document.body.classList.add('light-theme');
            document.body.classList.remove('dark-theme');
            themeToggle.checked = false;
        } else {
            document.body.classList.add('dark-theme');
            document.body.classList.remove('light-theme');
            themeToggle.checked = true;
        }
        
        // Обработчик переключения
        themeToggle.addEventListener('change', (e) => {
            if (e.target.checked) {
                document.body.classList.remove('light-theme');
                document.body.classList.add('dark-theme');
                localStorage.setItem('yanwit_theme', 'dark');
            } else {
                document.body.classList.remove('dark-theme');
                document.body.classList.add('light-theme');
                localStorage.setItem('yanwit_theme', 'light');
            }
        });
    }
});

// ========== Навигация ==========
function setupNavigation() {
    const navItems = document.querySelectorAll('.nav-item');
    
    navItems.forEach(item => {
        item.addEventListener('click', (e) => {
            e.preventDefault();
            const page = item.dataset.page;
            
            navItems.forEach(nav => nav.classList.remove('active'));
            item.classList.add('active');
            
            // Скрываем все страницы
            const homePage = document.getElementById('home-page');
            const profilePage = document.getElementById('profile-page');
            const searchPage = document.getElementById('search-page');
            const settingsPage = document.getElementById('settings-page');
            
            if (homePage) homePage.style.display = 'none';
            if (profilePage) profilePage.style.display = 'none';
            if (searchPage) searchPage.style.display = 'none';
            if (settingsPage) settingsPage.style.display = 'none';
            
            // Показываем выбранную и загружаем данные
            switch(page) {
                case 'home':
                    if (homePage) {
                        homePage.style.display = 'block';
                        if (typeof loadTimeline === 'function') loadTimeline();
                    }
                    break;
                case 'profile':
                    if (profilePage) {
                        profilePage.style.display = 'block';
                        if (typeof loadProfile === 'function') loadProfile();
                    }
                    break;
                case 'search':
                    if (searchPage) {
                        searchPage.style.display = 'block';
                        initSearch();
                    }
                    break;
                case 'settings':
                    if (settingsPage) settingsPage.style.display = 'block';
                    break;
            }
        });
    });
}

// ========== Поиск пользователей ==========
function initSearch() {
    const searchInput = document.getElementById('search-input');
    const resultsDiv = document.getElementById('search-results');
    
    if (!searchInput) {
        console.log('Search input not found');
        return;
    }
    
    // Удаляем старый обработчик
    const newSearchInput = searchInput.cloneNode(true);
    searchInput.parentNode.replaceChild(newSearchInput, searchInput);
    
    newSearchInput.addEventListener('input', async (e) => {
        const query = e.target.value.trim();
        
        if (query.length < 2) {
            if (resultsDiv) resultsDiv.innerHTML = '';
            return;
        }
        
        if (resultsDiv) resultsDiv.innerHTML = '<div class="loading">Поиск...</div>';
        
        try {
            const users = await searchUsers(query);
            renderSearchResults(users, resultsDiv);
        } catch (error) {
            console.error('Search error:', error);
            if (resultsDiv) resultsDiv.innerHTML = '<div class="empty-state">Ошибка поиска</div>';
        }
    });
}

function renderSearchResults(users, resultsDiv) {
    if (!resultsDiv) return;
    
    if (!users || users.length === 0) {
        resultsDiv.innerHTML = '<div class="empty-state">👤 Пользователи не найдены</div>';
        return;
    }
    
    // Получаем текущего пользователя
    const currentUserId = currentUser ? currentUser.id : null;
    
    resultsDiv.innerHTML = users.map(user => {
        // Не показываем кнопку для самого себя
        if (user.id === currentUserId) {
            return `
                <div class="search-result" data-user-id="${user.id}">
                    <div class="search-result-info">
                        <div class="search-result-username">@${escapeHtml(user.username)}</div>
                        <div class="search-result-bio">${escapeHtml(user.bio || 'Нет описания')}</div>
                    </div>
                    <span class="badge-self">Это вы</span>
                </div>
            `;
        }
        
        return `
            <div class="search-result" data-user-id="${user.id}">
                <div class="search-result-info">
                    <div class="search-result-username">@${escapeHtml(user.username)}</div>
                    <div class="search-result-bio">${escapeHtml(user.bio || 'Нет описания')}</div>
                </div>
                <button class="btn-follow" data-user-id="${user.id}" onclick="followUser(${user.id})">
                    Подписаться
                </button>
            </div>
        `;
    }).join('');
}

async function followUser(userId) {
    try {
        const response = await apiRequest(`/users/${userId}/follow`, {
            method: 'POST'
        });
        
        if (response.ok) {
            alert('✅ Подписка оформлена');
            
            // Обновляем поиск (чтобы кнопка изменилась или исчезла)
            initSearch();
            
            // Обновляем профиль, если открыт профиль текущего пользователя
            const currentPage = document.querySelector('.nav-item.active').dataset.page;
            if (currentPage === 'profile') {
                loadProfile(); // перезагружаем профиль с новыми цифрами
            }
        } else {
            const error = await response.json();
            alert('❌ Ошибка: ' + (error.error || 'Не удалось подписаться'));
        }
    } catch (error) {
        console.error('Follow error:', error);
        alert('❌ Ошибка подписки');
    }
}

// ========== Авторизация ==========
function setupAuthTabs() {
    const tabs = document.querySelectorAll('.auth-tab');
    if (!tabs.length) return;
    
    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            tabs.forEach(t => t.classList.remove('active'));
            tab.classList.add('active');
            
            const loginForm = document.getElementById('login-form');
            const registerForm = document.getElementById('register-form');
            
            if (tab.dataset.tab === 'login') {
                if (loginForm) loginForm.style.display = 'flex';
                if (registerForm) registerForm.style.display = 'none';
            } else {
                if (loginForm) loginForm.style.display = 'none';
                if (registerForm) registerForm.style.display = 'flex';
            }
        });
    });
}

function setupAuthForms() {
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = document.getElementById('login-username').value;
            const password = document.getElementById('login-password').value;
            const errorDiv = document.getElementById('auth-error');
            
            try {
                await login(username, password);
                showMainApp();
                if (typeof loadTimeline === 'function') loadTimeline();
            } catch (err) {
                if (errorDiv) errorDiv.textContent = err.message;
            }
        });
    }
    
    const registerForm = document.getElementById('register-form');
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = document.getElementById('register-username').value;
            const email = document.getElementById('register-email').value;
            const password = document.getElementById('register-password').value;
            const errorDiv = document.getElementById('register-error');
            
            try {
                await register(username, email, password);
                showMainApp();
                if (typeof loadTimeline === 'function') loadTimeline();
            } catch (err) {
                if (errorDiv) errorDiv.textContent = err.message;
            }
        });
    }
}

function showMainApp() {
    const authPage = document.getElementById('auth-page');
    const homePage = document.getElementById('home-page');
    if (authPage) authPage.style.display = 'none';
    if (homePage) homePage.style.display = 'block';
}

function showAuthPage() {
    const authPage = document.getElementById('auth-page');
    const homePage = document.getElementById('home-page');
    if (authPage) authPage.style.display = 'flex';
    if (homePage) homePage.style.display = 'none';
}

// ========== Профиль ==========
async function loadProfile() {
    const profileDiv = document.getElementById('profile-header');
    const statsDiv = document.getElementById('profile-stats');
    const tweetsDiv = document.getElementById('profile-tweets');
    
    if (!profileDiv) return;
    
    profileDiv.innerHTML = '<div class="loading">Загрузка профиля...</div>';
    if (statsDiv) statsDiv.innerHTML = '';
    if (tweetsDiv) tweetsDiv.innerHTML = '';
    
    try {
        const response = await apiRequest('/profile');
        const data = await response.json();
        const user = data.user || data;
        
        const avatarLetter = user.username ? user.username.charAt(0).toUpperCase() : '?';
        
        if (profileDiv) {
            profileDiv.innerHTML = `
                <div class="profile-banner"></div>
                <div class="profile-avatar-large">${avatarLetter}</div>
                <div class="profile-info">
                    <div class="profile-name">@${escapeHtml(user.username)}</div>
                    <div class="profile-bio">${escapeHtml(user.bio || 'Нет описания')}</div>
                    <div class="profile-joined">📅 Присоединился: ${new Date(user.created_at).toLocaleDateString()}</div>
                </div>
            `;
        }
        
        if (statsDiv) {
            // Загружаем актуальную статистику через отдельный API
            try {
                const statsResponse = await apiRequest(`/users/stats?user_id=${user.id}`);
                const stats = await statsResponse.json();
                
                statsDiv.innerHTML = `
                    <div class="stat">
                        <div class="stat-count">${stats.tweet_count || 0}</div>
                        <div class="stat-label">Твитов</div>
                    </div>
                    <div class="stat">
                        <div class="stat-count">${stats.followers_count || 0}</div>
                        <div class="stat-label">Подписчиков</div>
                    </div>
                    <div class="stat">
                        <div class="stat-count">${stats.following_count || 0}</div>
                        <div class="stat-label">Подписок</div>
                    </div>
                `;
            } catch (error) {
                console.error('Error loading stats:', error);
                statsDiv.innerHTML = `
                    <div class="stat">
                        <div class="stat-count">${user.tweet_count || 0}</div>
                        <div class="stat-label">Твитов</div>
                    </div>
                    <div class="stat">
                        <div class="stat-count">${user.followers_count || 0}</div>
                        <div class="stat-label">Подписчиков</div>
                    </div>
                    <div class="stat">
                        <div class="stat-count">${user.following_count || 0}</div>
                        <div class="stat-label">Подписок</div>
                    </div>
                `;
            }
        }
        
        if (tweetsDiv && typeof getTweetsByUsername === 'function') {
            await getTweetsByUsername(user.username, tweetsDiv);
        } else if (tweetsDiv) {
            tweetsDiv.innerHTML = '<div class="empty-state">Функция загрузки твитов временно недоступна</div>';
        }
        
    } catch (error) {
        console.error('Error loading profile:', error);
        if (profileDiv) {
            profileDiv.innerHTML = '<div class="empty-state">Ошибка загрузки профиля</div>';
        }
    }
}

// ========== Вспомогательные функции ==========
function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}