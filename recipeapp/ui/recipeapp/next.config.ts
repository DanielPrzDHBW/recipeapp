import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "export",
  images: {
    remotePatterns: [new URL('https://www.themealdb.com/images/**')],
    unoptimized: true,
  },
};

module.exports = nextConfig;
