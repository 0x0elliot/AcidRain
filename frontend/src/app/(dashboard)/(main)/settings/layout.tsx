export default function Layout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
      <><h1 className="text-lg font-semibold text-gray-900 sm:text-xl dark:text-gray-50">
      Settings
    </h1><div>{children}</div></>
  )
}
