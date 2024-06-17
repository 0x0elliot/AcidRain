"use client"

import { overviews } from "@/data/overview-data"
import { subDays, toDate } from "date-fns"
import React from "react"
import { DateRange } from "react-day-picker"
import { LineChart } from '@tremor/react';

const chartdata = [
  {
    date: 'Jun 24',
    Subscribers: 0,
  },
]



const overviewsDates = overviews.map((item) => toDate(item.date).getTime())
const maxDate = toDate(Math.max(...overviewsDates))

export default function Overview() {
  const [selectedDates, setSelectedDates] = React.useState<
    DateRange | undefined
  >({
    from: subDays(maxDate, 30),
    to: maxDate,
  })

  const customTooltip = (props) => {
    const { payload, active } = props;
    if (!active || !payload) return null;
    return (
      <div className="w-56 rounded-tremor-default border border-tremor-border bg-tremor-background p-2 text-tremor-default shadow-tremor-dropdown">
        {payload.map((category, idx) => (
          <div key={idx} className="flex flex-1 space-x-2.5 dark:text-gray-50 text-gray-900">
            <div
              className={`flex w-1 flex-col bg-${category.color}-500 rounded`}
            />
            <div className="space-y-1">
              <p className="text-tremor-content">{category.payload.date}</p>
              <p className="font-medium text-tremor-content-emphasis">
                {category.value} Active Subscribers
              </p>
            </div>
          </div>
        ))}
      </div>
    );
  };


  return (
    <>
      <section aria-labelledby="flows-title">
        <h1
          id="overall-title"
          className="scroll-mt-10 text-lg font-semibold text-gray-900 sm:text-xl dark:text-gray-50"
        >
          Your Newsletter Performance
        </h1>
      </section>

      {/* Add spacing and then use imports to render a chart of active subsscribers with dummy data */}
      <div className="scroll-mt-10" style={{ marginTop: '20px'}}>
        <>
          <h3 className="text-lg font-medium text-tremor-content-strong text-gray-900 dark:text-gray-50">
            Subscribers
          </h3>
          <div className="sticky top-16 z-20 flex items-center justify-between border-b border-gray-200 bg-white pb-4 pt-4 sm:pt-6 lg:top-0 lg:mx-0 text-gray-900 lg:px-0 lg:pt-8 dark:border-gray-800 dark:bg-gray-950">
            <LineChart
              className="mt-4 h-72 text-gray-900 dark:text-gray-50"
              // className="h-72"
              data={chartdata}
              index="date"
              categories={['Subscribers']}
              colors={['blue']}
              // yAxisWidth={30}
              customTooltip={customTooltip}
            />
          </div>
        </>
      </div>

      

    </>
  )
}
