
import Link from "next/link"
import styles from '@/styles/Home.module.css'

export default function ProductsList(){
    return (
        <div className ={styles.container}>
            <main className={styles.main}>
                <h2 className={styles.title}>動画リスト</h2>

                <ul>
<li>
   <Link href="/products/first">
            初動画
        </Link>
    </li>
    <li>
       <Link href="/products/sing">
            歌ってみた
        </Link>
    </li>
    <li>
       <Link href="/products/game">
            ゲーム実況
        </Link>
    </li>
</ul>

            </main>
        </div>                
     
    )}