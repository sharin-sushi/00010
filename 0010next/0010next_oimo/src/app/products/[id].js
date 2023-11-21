import Link from "next/link"
import styles from '@/styles/Home.module.css'
import{ useRouter} from "next/router";

// SSGの場合
// export async function getStaticProps({params}){
//     const req = await fetch(`http://localhost:3000/${params.id}.json`);
//     const data = await req.json();
//     console.log(data)
//     return{
//       props: {
//         product: data,
//       },
//     };
// }

// export async function getStaticPaths(){
//     const req = await fetch(`http://localhost:3000/products.json`);
//     const data = await req.json();

//     const paths = data.map((product) => {
//       return{
//         params:{
//             id: product,
//         },
//       };
//     });

//     return{
//       paths,
//       fallback: false,
//     };
// }

export async function getServerSideProps({params}){
    const req = await fetch(`http://localhost:3000/${params.id}.json`);
    const data = await req.json();
 
    return{
      props: {
        product: data,
      },
    };
}

 const Product = ({product}) => {
  const router = useRouter();
  const { id } = router.query;
  console.log(product)
  return (
    <div className ={styles.container}>
        <main className={styles.main}>
          <h1>{id}のページです</h1>
          <img src={product.image} width="50%" height="50%" />
          <p>{product.name}</p>
          <br />
          <Link href="/products">
            動画一覧へ
          </Link>
        </main>
    </div>    
  )    
};

export default Product;