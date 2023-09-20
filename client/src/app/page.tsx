import Image from "next/image";

async function getProducts() {
  const res = await fetch("http://localhost:8080/product");

  const products = await res.json();

  return products;
}

export default async function Home() {
  const products = await getProducts();

  console.log(products);

  return (
    <>
      <div>TEMP</div>

      <p>{JSON.stringify(products)}</p>
    </>
  );
}
