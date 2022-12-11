// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("prism-react-renderer/themes/github");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "Pensieve",
  tagline: "Smarter Home for all your Photos and Videos",
  url: "https://pensieve-docs.vercel.app",
  baseUrl: "/",
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/favicon.ico",
  organizationName: "prabhuomkar",
  projectName: "pensieve",
  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },
  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      }),
    ],
  ],
  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        title: "Pensieve",
        logo: {
          alt: "Pensieve Logo",
          src: "img/logo.svg",
        },
        items: [
          {
            type: "doc",
            docId: "dev-guide/introduction",
            position: "right",
            label: "Documentation",
          },
          {
            href: "https://github.com/prabhuomkar/pensieve",
            label: "GitHub",
            position: "right",
          },
        ],
      },
      footer: {
        style: "dark",
        copyright: `Copyright © ${new Date().getFullYear()} Pensieve. Built with ❤️ in India.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
  stylesheets: [
    "https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500;700&family=Fira+Code&display=swap",
  ],
};

module.exports = config;
