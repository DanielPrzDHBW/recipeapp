'use client';
import Image from "next/image";
import { useEffect, useState } from "react";

export default function Home() {
  function Sleep(milliseconds: number) {
    return new Promise(resolve => setTimeout(resolve, milliseconds));
  }


  const [data, setData] = useState();
  const [refetch, setRefetch] = useState(false);



  useEffect(() => {
    const fetchData = async () => {
      await fetch(`http://localhost:8080/api/newrecipes`)
      .then(response => response.json())
      .then((datac) => {
            console.log(datac);
            setData(datac);
         })
      .catch(error => {
            alert('Error fetching data: ' + error);
         });

    }

    fetchData();
  }, [refetch])
    
  return (
    <div className="font-sans grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
      <main className="flex flex-col gap-[32px] row-start-2 items-center">
        <h1 className="text-4xl sm:text-5xl font-extrabold text-center">
          RecipeApp
        </h1>
        <div className="flex gap-4 items-center flex-col sm:flex-row">
          <a
            //className="rounded-full border border-solid border-transparent transition-colors flex items-center justify-center bg-foreground text-background gap-2 hover:bg-[#383838] dark:hover:bg-[#ccc] font-medium text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 sm:w-auto"
            className="rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent font-medium text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 w-full sm:w-auto gap-2"
            onClick={() => setRefetch(!refetch)}
            target="_blank"
            rel="noopener noreferrer"
          >
            <Image
              className="dark:invert"
              src="/recipe.svg"
              alt="Recipe icon"
              width={20}
              height={20}
            />
            Generate Recipes
          </a>

          <a
            className="rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent font-medium text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 w-full sm:w-auto gap-2 "
            onClick={() => alert('Generate Recipes clicked!')}
            target="_blank"
            rel="noopener noreferrer"
          >
            <Image
              className="dark:invert"
              src="/shoppinglist.svg"
              alt="shoppinglist icon"
              width={20}
              height={20}
            />
            Shopping List
          </a>
        </div>
        <div className="flex gap-4 items-center flex-col sm:w-9/12">
          <Recipes datar={data} />
        </div>
      </main>
      <footer className="row-start-3 flex gap-[24px] flex-wrap items-center justify-center">
      </footer>
    </div>
  );
}

var recipes = [
  {
    idMeal: "",
    strMeal: "",
    strMealAlternate: "",
    strCategory: "",
    strArea: "",
    strInstructions: "",
    strMealThumb: "",
    strTags: "",
    strYoutube: "",
    strIngredient1: "",
    strIngredient2: "",
    strIngredient3: "",
    strIngredient4: "",
    strIngredient5: "",
    strIngredient6: "",
    strIngredient7: "",
    strIngredient8: "",
    strIngredient9: "",
    strIngredient10: "",
    strIngredient11: "",
    strIngredient12: "",
    strIngredient13: "",
    strIngredient14: "",
    strIngredient15: "",
    strIngredient16: "",
    strIngredient17: "",
    strIngredient18: "",
    strIngredient19: "",
    strIngredient20: "",
    strMeasure1: "",
    strMeasure2: "",
    strMeasure3: "",
    strMeasure4: "",
    strMeasure5: "",
    strMeasure6: "",
    strMeasure7: "",
    strMeasure8: "",
    strMeasure9: "",
    strMeasure10: "",
    strMeasure11: "",
    strMeasure12: "",
    strMeasure13: "",
    strMeasure14: "",
    strMeasure15: "",
    strMeasure16: "",
    strMeasure17: "",
    strMeasure18: "",
    strMeasure19: "",
    strMeasure20: "",
    strSource: "",
    strImageSource: "",
    strCreativeCommonsConfirmed: "",
    dateModified: ""
  }
]
function Recipes({ datar }: { datar: any }) {
  if (!datar) {
    return <div></div>;
  }
  recipes = []
  recipes = datar.recipe

  function handleClick(val: string) {

  }

  return <>{recipes.map(recipe =>
    <a
      className="rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent font-medium text-sm sm:text-base h-min sm:h-min px-4 sm:px-5 py-1 sm:py-2 w-full gap-2"
      key={recipe.strMeal}
      onClick={() => handleClick("Clicked")}
    >
      <div className="flex gap4 items-center flex-col sm:flex-row">
        <Image
          src={recipe.strMealThumb}
          alt="Image Food"
          width={20}
          height={20}
        />
        <h1 className="font-bold text-center">{recipe.strMeal}</h1>
      </div>
    </a>)}</>;
}