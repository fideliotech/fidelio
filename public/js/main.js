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

  const contactForm = document.getElementById('contactForm');
  if (contactForm) {
    contactForm.addEventListener('submit', async function(e) {
      e.preventDefault();
      const formData = new FormData(contactForm);
      const data = {};
      formData.forEach((value, key) => {
        data[key] = value;
      });
      try {
        const response = await fetch('/contact-request', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(data),
        });
        if (response.ok) {
          showSnackbar('Request sent successfully!');
          contactForm.reset();
        } else {
          showSnackbar('Failed to send request. Please try again.');
        }
      } catch (err) {
        showSnackbar('Error sending request.');
      } 
    });
  }

  let snackbarTimeout = null;
  function showSnackbar(message) {
    const snackbar = document.getElementById('snackBar');
    const snackText = snackbar.querySelector('.snack_text');
    
    if (snackbarTimeout) {
      clearTimeout(snackbarTimeout);
      
      if (snackbar.classList.contains('show')) {
        snackbar.classList.remove('show');
          
        setTimeout(() => {
          showNewSnackbar();
        }, 50);
      } else {
        showNewSnackbar();
      }
    } else {
      showNewSnackbar();
    }
    
    function showNewSnackbar() {
      snackText.textContent = message;
      snackbar.classList.add('show');
      
      snackbarTimeout = setTimeout(() => {
        snackbar.classList.remove('show');
        snackbarTimeout = null;
      }, 4000);
    }
  }

  const languageSwitcher = document.getElementById('language-switcher');
  const drawerLanguageSwitcher = document.getElementById('drawer-language-switcher');
  const currentLang = localStorage.getItem('language') || 'es';

  function updateLanguage(lang) {
    if (lang !== 'en' && lang !== 'es') return;
    
    localStorage.setItem('language', lang);
    document.documentElement.lang = lang;
    
    const newLang = lang === 'en' ? 'es' : 'en';
    const switchText = translations[newLang]['language.switch'];
    
    // Update both language switchers
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
    
    // Update all translatable elements
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

  // Add click handlers for both language switchers
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

  function scrollToSection(element) {
    const headerOffset = 110;
    const elementPosition = element.getBoundingClientRect().top;
    const offsetPosition = elementPosition + window.pageYOffset - headerOffset;
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;

    if (prefersReducedMotion) {
      window.scrollTo({
        top: offsetPosition,
        behavior: 'auto'
      });
    } else {
      const startPosition = window.scrollY;
      const distance = offsetPosition - startPosition;
      const duration = 1200;
      let startTime = null;

      function easeInOutCubic(t) {
        return t < 0.5 ? 4 * t * t * t : (t - 1) * (2 * t - 2) * (2 * t - 2) + 1;
      }

      function animation(currentTime) {
        if (startTime === null) startTime = currentTime;
        const timeElapsed = currentTime - startTime;
        const progress = Math.min(timeElapsed / duration, 1);
        const ease = easeInOutCubic(progress);
        
        window.scrollTo(0, startPosition + distance * ease);
        
        if (timeElapsed < duration) {
          requestAnimationFrame(animation);
        }
      }

      requestAnimationFrame(animation);
    }

    document.documentElement.style.overflow = '';
    body.style.overflow = '';
  }

  document.querySelectorAll('.nav_container a, .drawer_content a, .hero_cta a').forEach(link => {
    link.addEventListener('click', handleNavigationClick);
  });
})