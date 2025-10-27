'use client';
import Image from "next/image";
import { MouseEventHandler, useEffect, useState, useRef } from "react";
import Popup from 'reactjs-popup';
import useSWR, { useSWRConfig } from 'swr'
import useSWRMutation from 'swr/mutation'
import React from "react";

const fetcher = (url: string) => fetch(url).then((res) => res.json());


export default function Home() {

  const { mutate } = useSWRConfig()
  const { data, trigger:newrecipe } = useSWRMutation('http://localhost:8080/api/newrecipes', fetcher)
  const [showNew, setShowNew] = useState(false);

  return (
    <div className="font-sans grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
      <main className="flex flex-col gap-[32px] row-start-2 items-center">
        <h1 className="text-4xl sm:text-5xl font-extrabold text-center">
          RecipeApp
        </h1>
        <div className="flex gap-4 items-center flex-col sm:flex-row">
          <a
            className="rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent font-medium text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 w-full sm:w-auto gap-2"
            onClick={async () => {
              await newrecipe();
              setShowNew(true);
            }}
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
          <RecipesComp extData={showNew && data ? data : undefined} />
        </div>
      </main>
      <footer className="row-start-3 flex gap-[24px] flex-wrap items-center justify-center">
      </footer>
    </div>
  );
}

type Recipe =
  {
    idMeal: string,
    strMeal: string,
    strMealAlternate: string,
    strCategory: string,
    strArea: string,
    strInstructions: string,
    strMealThumb: string,
    strTags: string,
    strYoutube: string,
    strIngredient1: string,
    strIngredient2: string,
    strIngredient3: string,
    strIngredient4: string,
    strIngredient5: string,
    strIngredient6: string,
    strIngredient7: string,
    strIngredient8: string,
    strIngredient9: string,
    strIngredient10: string,
    strIngredient11: string,
    strIngredient12: string,
    strIngredient13: string,
    strIngredient14: string,
    strIngredient15: string,
    strIngredient16: string,
    strIngredient17: string,
    strIngredient18: string,
    strIngredient19: string,
    strIngredient20: string,
    strMeasure1: string,
    strMeasure2: string,
    strMeasure3: string,
    strMeasure4: string,
    strMeasure5: string,
    strMeasure6: string,
    strMeasure7: string,
    strMeasure8: string,
    strMeasure9: string,
    strMeasure10: string,
    strMeasure11: string,
    strMeasure12: string,
    strMeasure13: string,
    strMeasure14: string,
    strMeasure15: string,
    strMeasure16: string,
    strMeasure17: string,
    strMeasure18: string,
    strMeasure19: string,
    strMeasure20: string,
    strSource: string,
    strImageSource: string,
    strCreativeCommonsConfirmed: string,
    dateModified: string
  }

type Recipes = {
  recipe: Recipe[]
}

function RecipePopupContent({ recipe, close }: { recipe: Recipe, close: () => void }) {
  const contentRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (contentRef.current) {
      contentRef.current.scrollTop = 0;
    }
  }, [recipe]); // Runs when recipe changes (i.e., popup opens for a new recipe)

  return (
    <div
      ref={contentRef}
      className="overflow-y-auto bg-gray-800 p-8 rounded shadow-lg flex flex-col items-center max-h-[90vh] max-w-[80vw]"
    >
      <h1 className="text-2xl font-bold mb-4">{recipe.strMeal}</h1>
      <div>
        <p className="mb-4 whitespace-pre-wrap text-xs">ID:{recipe.idMeal} | Category: {recipe.strCategory}</p>
      </div>
      <Image
        src={recipe.strMealThumb}
        alt="Image Food"
        width={300}
        height={300}
      />
      <div>
        <h2 className="text-xl font-semibold mb-2 text-center">Ingredients</h2>
        <ul className="text-center">
          {[...Array(20)].map((_, i) => {
            const ingredient = recipe[`strIngredient${i + 1}` as keyof Recipe];
            const measure = recipe[`strMeasure${i + 1}` as keyof Recipe];
            return ingredient ? (
              <li key={i}>{ingredient} - {measure}</li>
            ) : null;
          })}
        </ul>
      </div>
      <div>
        <h2 className="text-xl font-semibold mb-2 text-center">Instructions</h2>
        <p className="mb-4 whitespace-pre-wrap text-m text-center">{recipe.strInstructions}</p>
      </div>
      <a className="text-xs" href={recipe.strYoutube}>Youtube</a>
      <button
        className="mt-4 px-4 py-2 bg-gray-200 text-gray-500 rounded"
        onClick={close}
      >
        Close
      </button>
    </div>
  );
}

function RecipesComp({ extData }: { extData?: Recipes }) {
  // Only fetch from SWR if extData is not present
  const { data, error, isLoading } = !extData
    ? useSWR<Recipes>('http://localhost:8080/api/recipes', fetcher)
    : { data: extData, error: undefined, isLoading: false };

  if (error) return <div>Failed to load</div>
  if (isLoading) return <div>Loading...</div>
  if (!data) return null;

  return (
    <>
      {data.recipe.map((recipe: Recipe) =>
        <Popup
          key={recipe.strMeal}
          trigger={
            <a
              className="rounded-full border border-solid border-black/[.08] dark:border-white/[.145] transition-colors flex items-center justify-center hover:bg-[#f2f2f2] dark:hover:bg-[#1a1a1a] hover:border-transparent font-medium text-sm sm:text-base h-min sm:h-min px-4 sm:px-5 py-1 sm:py-2 w-full gap-2"
            >
              <div className="flex gap4 items-center flex-col sm:flex-row">
                <Image
                  src={recipe.strMealThumb}
                  alt="Image Food"
                  width={50}
                  height={50}
                />
                <h1 className="font-bold text-center">{recipe.strMeal}</h1>
              </div>
            </a>
          }
          modal
          nested
          contentStyle={{ padding: 0, border: "none", background: "none" }}
        >
          {((close: () => void) => (
            <RecipePopupContent recipe={recipe} close={close} />
          )) as any}
        </Popup>
      )}
    </>
  );
}
