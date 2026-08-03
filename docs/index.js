const themes = {
  dark: {
    bg: "#0d1117",
    text: "#c9d1d9",
    subtext: "#8b949e",
    barbg: "#21262d",
    c1: "#39d353",
    c2: "#2ea043",
    c3: "#238636",
    c4: "#196c2e",
  },
  light: {
    bg: "#FDF6E3",
    text: "#657B83",
    subtext: "#93A1A1",
    barbg: "#EEE8D5",
    c1: "#268BD2",
    c2: "#859900",
    c3: "#B58900",
    c4: "#D33682",
  },
  dracula: {
    bg: "#282a36",
    text: "#f8f8f2",
    subtext: "#6272a4",
    barbg: "#44475a",
    c1: "#bd93f9",
    c2: "#50fa7b",
    c3: "#ff79c6",
    c4: "#8be9fd",
  },
  nord: {
    bg: "#2e3440",
    text: "#d8dee9",
    subtext: "#4c566a",
    barbg: "#3b4252",
    c1: "#88c0d0",
    c2: "#a3be8c",
    c3: "#81a1c1",
    c4: "#bf616a",
  },
  gruvbox: {
    bg: "#282828",
    text: "#ebdbb2",
    subtext: "#928374",
    barbg: "#3c3836",
    c1: "#fabd2f",
    c2: "#b8bb26",
    c3: "#fe8019",
    c4: "#fb4934",
  },
  monokai: {
    bg: "#272822",
    text: "#f8f8f2",
    subtext: "#75715e",
    barbg: "#3e3d32",
    c1: "#a6e22e",
    c2: "#f92672",
    c3: "#66d9ef",
    c4: "#fd971f",
  },
  cyberpunk: {
    bg: "#0b0e14",
    text: "#e0f7fa",
    subtext: "#ff0055",
    barbg: "#181c25",
    c1: "#00ff9f",
    c2: "#f6019d",
    c3: "#03a9f4",
    c4: "#ffeb3b",
  },
  tokyonight: {
    bg: "#1a1b26",
    text: "#c0caf5",
    subtext: "#565f89",
    barbg: "#24283b",
    c1: "#7aa2f7",
    c2: "#9ece6a",
    c3: "#e0af68",
    c4: "#f7768e",
  },
  everforest: {
    bg: "#2b3339",
    text: "#d3c6aa",
    subtext: "#7a8478",
    barbg: "#374145",
    c1: "#a7c080",
    c2: "#7fbbb3",
    c3: "#dbbc7f",
    c4: "#e67e80",
  },
  iceberg: {
    bg: "#161821",
    text: "#d2d4de",
    subtext: "#6b7089",
    barbg: "#1e2132",
    c1: "#84a0c6",
    c2: "#a093c7",
    c3: "#89b8c2",
    c4: "#e27878",
  },
  sunset: {
    bg: "#1f1d2b",
    text: "#f8f8f2",
    subtext: "#a599e9",
    barbg: "#2a273f",
    c1: "#ff9e64",
    c2: "#ffd580",
    c3: "#ff6b6b",
    c4: "#c678dd",
  },
  deepocean: {
    bg: "#0f172a",
    text: "#e2e8f0",
    subtext: "#64748b",
    barbg: "#1e293b",
    c1: "#38bdf8",
    c2: "#22c55e",
    c3: "#f59e0b",
    c4: "#ef4444",
  },
  midnightpurple: {
    bg: "#1b1325",
    text: "#e9d8fd",
    subtext: "#9f7aea",
    barbg: "#2d1b3f",
    c1: "#c084fc",
    c2: "#60a5fa",
    c3: "#34d399",
    c4: "#f472b6",
  },
  catppuccin: {
    bg: "#1e1e2e",
    text: "#cdd6f4",
    subtext: "#6c7086",
    barbg: "#313244",
    c1: "#89b4fa",
    c2: "#a6e3a1",
    c3: "#f9e2af",
    c4: "#f38ba8",
  },
  solarized: {
    bg: "#002b36",
    text: "#93a1a1",
    subtext: "#586e75",
    barbg: "#073642",
    c1: "#268bd2",
    c2: "#859900",
    c3: "#b58900",
    c4: "#dc322f",
  },
  onedark: {
    bg: "#282c34",
    text: "#abb2bf",
    subtext: "#5c6370",
    barbg: "#3a3f4b",
    c1: "#61afef",
    c2: "#98c379",
    c3: "#e5c07b",
    c4: "#e06c75",
  },
  material: {
    bg: "#263238",
    text: "#eeffff",
    subtext: "#546e7a",
    barbg: "#37474f",
    c1: "#82aaff",
    c2: "#c3e88d",
    c3: "#ffcb6b",
    c4: "#f07178",
  },
  synthwave: {
    bg: "#241b2f",
    text: "#f8f8f2",
    subtext: "#ff7edb",
    barbg: "#34294f",
    c1: "#36f9f6",
    c2: "#72f1b8",
    c3: "#fede5d",
    c4: "#ff5c8a",
  },
  rosepine: {
    bg: "#191724",
    text: "#e0def4",
    subtext: "#908caa",
    barbg: "#26233a",
    c1: "#eb6f92",
    c2: "#f6c177",
    c3: "#ebbcba",
    c4: "#31748f",
  },
};

const inputs = {
  bg: document.getElementById("bg"),
  text: document.getElementById("text"),
  subtext: document.getElementById("subtext"),
  barbg: document.getElementById("barbg"),
  c1: document.getElementById("c1"),
  c2: document.getElementById("c2"),
  c3: document.getElementById("c3"),
  c4: document.getElementById("c4"),
};

const themeSelector = document.getElementById("theme");
const output = document.getElementById("output");
const copyBtn = document.getElementById("copyBtn");
const copyBtnText = document.getElementById("copyBtnText");

function updateDummyCard(colors) {
  const card = document.getElementById("dummyCard");
  card.style.backgroundColor = colors.bg;

  document
    .querySelectorAll(".card-bg-match")
    .forEach((e) => (e.style.backgroundColor = colors.bg));

  document.getElementById("cardTitle").style.color = colors.c1;

  document
    .querySelectorAll(".card-subtext")
    .forEach((e) => (e.style.color = colors.subtext));
  document
    .querySelectorAll(".card-text")
    .forEach((e) => (e.style.color = colors.text));

  document
    .querySelectorAll(".var-c1")
    .forEach((e) => (e.style.color = colors.c1));
  document
    .querySelectorAll(".var-c2")
    .forEach((e) => (e.style.color = colors.c2));
  document
    .querySelectorAll(".var-c3")
    .forEach((e) => (e.style.color = colors.c3));
  document
    .querySelectorAll(".var-c4")
    .forEach((e) => (e.style.color = colors.c4));

  document
    .querySelectorAll(".card-bar-bg")
    .forEach((e) => (e.style.backgroundColor = colors.barbg));

  document
    .querySelectorAll(".card-c1")
    .forEach((e) => (e.style.backgroundColor = colors.c1));
  document
    .querySelectorAll(".card-c2")
    .forEach((e) => (e.style.backgroundColor = colors.c2));
  document
    .querySelectorAll(".card-c3")
    .forEach((e) => (e.style.backgroundColor = colors.c3));
  document
    .querySelectorAll(".card-c4")
    .forEach((e) => (e.style.backgroundColor = colors.c4));
}

function generateCommand(colors) {
  return `./taka-report -days=7 -bg "${colors.bg}" -text "${colors.text}" -subtext "${colors.subtext}" -bar-bg "${colors.barbg}" -c1 "${colors.c1}" -c2 "${colors.c2}" -c3 "${colors.c3}" -c4 "${colors.c4}"`;
}

function applyColors(colors) {
  Object.keys(inputs).forEach((key) => {
    inputs[key].value = colors[key];
  });
  updateDummyCard(colors);
  output.value = generateCommand(colors);
}

window.addEventListener("load", () => {
  setTimeout(() => {
    document.querySelectorAll(".progress-bar-fill").forEach((bar) => {
      const targetWidth = bar.getAttribute("data-width");
      if (targetWidth) {
        bar.style.width = targetWidth;
      }
    });
  }, 300);
});

themeSelector.addEventListener("change", (e) => {
  const selectedTheme = themes[e.target.value];
  if (selectedTheme) {
    applyColors(selectedTheme);
  }
});

Object.values(inputs).forEach((input) => {
  input.addEventListener("input", () => {
    const currentColors = {
      bg: inputs.bg.value,
      text: inputs.text.value,
      subtext: inputs.subtext.value,
      barbg: inputs.barbg.value,
      c1: inputs.c1.value,
      c2: inputs.c2.value,
      c3: inputs.c3.value,
      c4: inputs.c4.value,
    };
    updateDummyCard(currentColors);
    output.value = generateCommand(currentColors);
  });
});

copyBtn.addEventListener("click", () => {
  navigator.clipboard.writeText(output.value);

  const orig = copyBtnText.innerText;
  copyBtnText.textContent = "Copied! ✨";
  copyBtn.classList.remove("bg-[#238636]", "hover:bg-[#2ea043]");
  copyBtn.classList.add("bg-emerald-600", "hover:bg-emerald-500");
  setTimeout(() => {
    copyBtnText.innerText = orig;
    copyBtn.classList.add("bg-[#238636]", "hover:bg-[#2ea043]");
    copyBtn.classList.remove("bg-emerald-600", "hover:bg-emerald-500");
  }, 2000);

  // Toast
  const toast = document.getElementById("toast");
  toast.classList.add("show");
  setTimeout(() => toast.classList.remove("show"), 3000);
});

const defaultTheme = {
  bg: "#0d1117",
  text: "#00FF00",
  subtext: "#008800",
  barbg: "#111111",
  c1: "#00FF00",
  c2: "#00DD00",
  c3: "#00AA00",
  c4: "#005500",
};
applyColors(defaultTheme);

// Scroll Reveal Animation
const observerOptions = {
  threshold: 0.1,
  rootMargin: "0px 0px -50px 0px",
};

const revealObserver = new IntersectionObserver((entries) => {
  entries.forEach((entry) => {
    if (entry.isIntersecting) {
      entry.target.classList.add("revealed");
      revealObserver.unobserve(entry.target);
    }
  });
}, observerOptions);

// Apply observer to all scroll-reveal elements
document
  .querySelectorAll(
    ".scroll-reveal, .scroll-reveal-left, .scroll-reveal-right, .stagger-item",
  )
  .forEach((el) => {
    revealObserver.observe(el);
  });

// Apply to feature cards and table rows
document.querySelectorAll(".card, table tbody tr").forEach((el) => {
  el.classList.add("stagger-item");
  revealObserver.observe(el);
});

// Apply to section headings and paragraphs
document.querySelectorAll("section > h2, section > p").forEach((el) => {
  el.classList.add("scroll-reveal");
  revealObserver.observe(el);
});

// Smooth scroll behavior
document.documentElement.style.scrollBehavior = "smooth";

// Card Mouse Tracking Animation
const cards = document.querySelectorAll(".card");

cards.forEach((card) => {
  card.addEventListener("mousemove", (e) => {
    const rect = card.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    // Calculate percentage position
    const xPercent = (x / rect.width) * 100;
    const yPercent = (y / rect.height) * 100;

    // Set CSS variables for gradient
    card.style.setProperty("--mouse-x", `${xPercent}%`);
    card.style.setProperty("--mouse-y", `${yPercent}%`);

    // Calculate rotation based on position
    const rotateX = (yPercent - 50) * 0.3;
    const rotateY = (xPercent - 50) * 0.3;

    card.style.transform = `translateY(-12px) scale(1.05) rotateX(${rotateX}deg) rotateY(${rotateY}deg)`;
  });

  card.addEventListener("mouseleave", () => {
    card.style.transform = "";
  });
});

// Button Mouse Tracking Animation
const buttons = document.querySelectorAll(".btn");

buttons.forEach((btn) => {
  btn.addEventListener("mousemove", (e) => {
    const rect = btn.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    const xPercent = (x / rect.width) * 100;
    const yPercent = (y / rect.height) * 100;

    btn.style.setProperty("--mouse-x", `${xPercent}%`);
    btn.style.setProperty("--mouse-y", `${yPercent}%`);
  });

  btn.addEventListener("mouseleave", () => {
    btn.style.setProperty("--mouse-x", "50%");
    btn.style.setProperty("--mouse-y", "50%");
  });
});

AsciinemaPlayer.create("demo.cast", document.getElementById("demo"), {
  autoplay: true,
  loop: true,
  theme: "tokyo-night",
});

async function fetchGitHubStats() {
  const repo = "Rtarun3606k/TakaTime";

  try {
    // REPO DATA

    const repoRes = await fetch(`https://api.github.com/repos/${repo}`);

    const repoData = await repoRes.json();

    // RELEASE DOWNLOADS

    const releaseRes = await fetch(
      `https://api.github.com/repos/${repo}/releases`,
    );

    const releases = await releaseRes.json();

    let totalDownloads = 0;

    releases.forEach((release) => {
      release.assets.forEach((asset) => {
        totalDownloads += asset.download_count;
      });
    });

    animateCounter("stars-count", repoData.stargazers_count);

    animateCounter("downloads-count", totalDownloads);
  } catch (err) {
    console.error(err);
  }
}

function animateCounter(id, target) {
  const element = document.getElementById(id);

  let current = 0;

  const increment = target / 80;

  const timer = setInterval(() => {
    current += increment;

    if (current >= target) {
      current = target;

      clearInterval(timer);
    }

    element.innerText = Math.floor(current).toLocaleString() + "+";
  }, 20);
}

fetchGitHubStats();
