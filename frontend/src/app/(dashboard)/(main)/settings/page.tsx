import { Button } from "@/components/ui/button"

export default function Settings() {
  return (
    <>
      {/* Add import subscribers with csv or copy paste */}
      <section aria-labelledby="settings-import-subscribers" style={{ margin: '10px'}}>
        <h2
          id="settings-import-subscribers"
          className="scroll-mt-10 text-lg font-semibold text-gray-500 sm:text-xl dark:text-gray-50"
        >
          Import Subscribers
        </h2>

        <div className="flex flex-col space-y-2">
          <Button
            className="w-full"
          >
            Import subscribers from CSV
          </Button>
        </div>

        <div className="mt-4">
          <p className="text-gray-500 dark:text-gray-400">
            Or import subscribers from a CSV file or copy paste them into the text
            area below (Separate emails with commas).
          </p>

          <textarea
            className="w-full h-48 mt-2 p-2 border border-gray-300 dark:border-gray-700 rounded-md"
            placeholder="email1@gmail.com, email2@gmail.com.."
          ></textarea>
        </div>

        <div className="mt-4">
          <Button
            className="w-full"
          >
            Import subscribers
          </Button>
        </div>

        <div className="mt-4">
          <p className="text-gray-500 dark:text-gray-400">
            Subscribers will be added to your existing list.
          </p>
        </div>


      </section>
    </>
  )
}
