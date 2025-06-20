document.addEventListener("DOMContentLoaded", () => {
	const body = document.body;
	const navContainer = document.querySelector('.nav_container');
  const drawerToggle = document.getElementById('drawerToggle');
	const drawerLink = document.querySelector('.drawer_link');
	const whatsappContainer = document.querySelector('.whatsapp_container');

	function getScrollbarWidth() {
		return window.innerWidth - document.documentElement.clientWidth;
	}

	drawerToggle.addEventListener('change', () => {
		if (drawerToggle.checked) {
			const scrollbarWidth = getScrollbarWidth();
			body.style.overflow = 'hidden';
			body.style.paddingRight = `${scrollbarWidth}px`;
			navContainer.style.marginRight = `${scrollbarWidth}px`;
			drawerLink.style.marginRight = `${scrollbarWidth}px`;
			whatsappContainer.style.paddingRight = `calc(var(--container-padx) + ${scrollbarWidth}px)`;
		} else {
			body.style.overflow = '';
			body.style.paddingRight = '';
			navContainer.style.marginRight = '';
			drawerLink.style.marginRight = '';
			whatsappContainer.style.paddingRight = 'var(--container-padx)';
			document.documentElement.style.overflow = '';
			body.style.overflow = '';
		}
	});

  const languageSwitcher = document.getElementById('language-switcher');
  const drawerLanguageSwitcher = document.getElementById('drawer-language-switcher');
  const currentLang = localStorage.getItem('language') || 'es';

  function updateLanguage(lang) {
    if (lang !== 'en' && lang !== 'es') return;
    
    localStorage.setItem('language', lang);
    document.documentElement.lang = lang;
    
    const newLang = lang === 'en' ? 'es' : 'en';
    const switchText = translations[newLang]['language.switch'];
    
    [languageSwitcher, drawerLanguageSwitcher].forEach(switcher => {
      if (switcher) {
        switcher.setAttribute('data-lang', newLang);
        const textSpan = switcher.querySelector('span:not(.material-icons)');
        if (textSpan) {
          textSpan.textContent = switchText;
        } else {
          switcher.textContent = switchText;
        }
      }
    });
    
    document.querySelectorAll('[data-i18n]').forEach(element => {
      const key = element.getAttribute('data-i18n');
      if (translations[lang] && translations[lang][key]) {
        if (element.tagName === 'INPUT' && element.hasAttribute('data-i18n-value')) {
          element.value = translations[lang][key];
          element.textContent = translations[lang][key];
        } else {
          element.textContent = translations[lang][key];
        }
      }
    });
    
    document.querySelectorAll('[data-i18n-placeholder]').forEach(element => {
      const key = element.getAttribute('data-i18n-placeholder');
      if (translations[lang] && translations[lang][key]) {
        element.placeholder = translations[lang][key];
      }
    });
  }
  updateLanguage(currentLang);

  [languageSwitcher, drawerLanguageSwitcher].forEach(switcher => {
    if (switcher) {
      switcher.addEventListener('click', () => {
        const newLang = switcher.getAttribute('data-lang');
        updateLanguage(newLang);
      });
    }
  });

  function handleNavigationClick(e) {
    const target = e.target.closest('a');
    if (!target || !target.hash) return;

    e.preventDefault();
    const targetId = target.hash.substring(1);
    const targetElement = document.getElementById(targetId);
    
    if (!targetElement) return;

    if (drawerToggle.checked) {
      drawerToggle.checked = false;
      body.style.overflow = '';
      body.style.paddingRight = '';
      navContainer.style.marginRight = '';
      drawerLink.style.marginRight = '';
      setTimeout(() => {
        scrollToSection(targetElement);
      }, 300);
    } else {
      scrollToSection(targetElement);
    }
  }

  document.querySelectorAll('.nav_container a, .drawer_content a, .hero_cta a').forEach(link => {
    link.addEventListener('click', handleNavigationClick);
  });
})