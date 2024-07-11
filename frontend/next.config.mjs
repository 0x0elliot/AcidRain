/** @type {import('next').NextConfig} */

const nextConfig = {
  redirects: async () => {
    return [
      {
        source: "/",
        destination: "/automations",
        permanent: true,
      },
    ];
  },
};

export default nextConfig;
