const lightCodeTheme = require("prism-react-renderer/themes/github");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

const config = {
  title: "Carousel",
  tagline: "Smarter Home for all your Photos and Videos",
  url: "https://carousel.omkar.xyz",
  baseUrl: "/",
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/favicon.ico",
  organizationName: "prabhuomkar",
  projectName: "carousel",
  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },
  presets: [
    [
      "classic",
      ({
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      }),
    ],
    [
      'redocusaurus',
      {
        specs: [
          {
            spec: 'swagger.yaml',
            route: '/api/',
          },
        ],
      },
    ],
  ],
  themeConfig:
    ({
      announcementBar: {
        id: 'wip',
        content:
          'Currently under development, check out <a href="/docs/dev-guide/contribution">Contributing Guide</a>',
        backgroundColor: '#8390e7',
        textColor: '#ffffff',
        isCloseable: false,
      },
      navbar: {
        title: "Carousel",
        logo: {
          alt: "Carousel Logo",
          src: "img/logo.png",
          srcDark: "img/logo-white.png"
        },
        items: [
          {
            type: "doc",
            docId: "dev-guide/introduction",
            label: "Docs",
            position: "left",
          },
          {
            href: "/api/",
            label: "API",
            position: "left",
          },
          {
            type: "localeDropdown",
            position: "right",
          },
          {
            href: "https://github.com/prabhuomkar/carousel",
            position: "right",
            className: "header-github-link",
            "aria-label": "GitHub repository",
          },
        ],
      },
      footer: {
        style: "dark",
        copyright: `Copyright © ${new Date().getFullYear()} Carousel. Built with ❤️ in India.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
  stylesheets: [
    "https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500;700&family=Fira+Code&display=swap",
  ]
};

module.exports = config;
