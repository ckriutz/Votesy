'use client'

import { useState, useEffect } from 'react'

export default function Header() {
    const [currentQuestion, setCurrentQuestion] = useState(null)
    const [AllQuestions, setAllQuestions] = useState(null)
    const [isLoading, setLoading] = useState(true)

    useEffect(() => {
        fetch('http://localhost:10000/question/current')
          .then((res) => res.json())
          .then((data) => {
            setCurrentQuestion(data)
            setLoading(false)
          })
          fetch('http://localhost:10000/questions')
            .then((res) => res.json())
            .then((data) => {
              setAllQuestions(data)
              setLoading(false)
            })
      }, [])

    if(isLoading) return (<div>Loading...</div>)

    return (
        <div className="w-full overflow-x-auto shadow-md sm:rounded-lg m-2">
            <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
                <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                    <tr>
                        <th scope="col" className="px-6 py-3">
                            Question
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Answer 1
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Votes
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Answer 2
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Votes
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Creted Date
                        </th>
                    </tr>
                </thead>
                <tbody>
                    <tr className="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700">
                        <th scope="row" className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                            {AllQuestions[0].text}
                        </th>
                        <td className="px-6 py-4">
                            {AllQuestions[0].answer1Text}
                        </td>
                        <td className="px-6 py-4">
                            30
                        </td>
                        <td className="px-6 py-4">
                            {AllQuestions[0].answer2Text}
                        </td>
                        <td className="px-6 py-4">
                            45
                        </td>
                        <td className="px-6 py-4">
                        <button data-modal-target="crud-modal" data-modal-toggle="crud-modal" className="block text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800" type="button" >Toggle modal</button>
                        </td>
                    </tr>
                    <tr className="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700">
                        <th scope="row" className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                            Microsoft Surface Pro
                        </th>
                        <td className="px-6 py-4">
                            White
                        </td>
                        <td className="px-6 py-4">
                            Laptop PC
                        </td>
                        <td className="px-6 py-4">
                            $1999
                        </td>
                        <td className="px-6 py-4">
                            <a href="#" className="font-medium text-blue-600 dark:text-blue-500 hover:underline">Edit</a>
                        </td>
                    </tr>
                    <tr className="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700">
                        <th scope="row" className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                            Magic Mouse 2
                        </th>
                        <td className="px-6 py-4">
                            Black
                        </td>
                        <td className="px-6 py-4">
                            Accessories
                        </td>
                        <td className="px-6 py-4">
                            $99
                        </td>
                        <td className="px-6 py-4">
                            <a href="#" className="font-medium text-blue-600 dark:text-blue-500 hover:underline">Edit</a>
                        </td>
                    </tr>
                    <tr className="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700">
                        <th scope="row" className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                            Google Pixel Phone
                        </th>
                        <td className="px-6 py-4">
                            Gray
                        </td>
                        <td className="px-6 py-4">
                            Phone
                        </td>
                        <td className="px-6 py-4">
                            $799
                        </td>
                        <td className="px-6 py-4">
                            <a href="#" className="font-medium text-blue-600 dark:text-blue-500 hover:underline">Edit</a>
                        </td>
                    </tr>
                    <tr>
                        <th scope="row" className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                            Apple Watch 5
                        </th>
                        <td className="px-6 py-4">
                            Red
                        </td>
                        <td className="px-6 py-4">
                            Wearables
                        </td>
                        <td className="px-6 py-4">
                            $999
                        </td>
                        <td className="px-6 py-4">
                            <a href="#" className="font-medium text-blue-600 dark:text-blue-500 hover:underline">Edit</a>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    )
}