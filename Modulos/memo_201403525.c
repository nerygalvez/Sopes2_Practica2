/**
*
*	Creando el modulo de memoria
*
*/


#include <linux/module.h> 

#include <linux/kernel.h>

#include <linux/init.h>
#include <linux/list.h>
#include <linux/types.h>
#include <linux/slab.h>
#include <linux/sched.h>
#include <linux/string.h>
#include <linux/fs.h>
#include <linux/seq_file.h>
#include <linux/proc_fs.h>
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>

#define FileProc "201403525_memo"
#define Carne "201403525"
#define Nombre "Nery Gonzalo Galvez Gomez"
#define SO "Ubuntu 18.04.4 LTS"
#define Curso "Sistemas Operativos 2"



struct sysinfo i;

static int proc_llenar_archivo(struct seq_file *m, void *v) {
	#define S(x) ((x) << (PAGE_SHIFT -10))
	si_meminfo(&i);

        //Aquí voy a escribir el archivo con formato json
        seq_printf(m, "{\n\"Total\" : %8lu \n", S(i.totalram)/1000);
        seq_printf(m, ", \"Consumida\" : %8lu \n", (S(i.totalram) - S(i.freeram))/1000);
        seq_printf(m, ", \"Porcentaje\" : %8lu \n}",100 - (S(i.freeram)*100/S(i.totalram)));      
        
        return 0;
}

static int proc_al_abrir_archivo(struct inode *inode, struct  file *file) {
  return single_open(file, proc_llenar_archivo, NULL);
}

static struct file_operations myops =
{
        .owner = THIS_MODULE,
        .open = proc_al_abrir_archivo,
        .read = seq_read,
        .llseek = seq_lseek,
        .release = single_release,
};



/**
*	Defino que es lo que se va a hacer al cargar el modulo
*/
static int iniciar(void)
{
	proc_create(FileProc,0,NULL,&myops);
    	printk(KERN_INFO "Carne: %s\n", Carne);

        /*
         * Si no se devuelve 0 significa que initmodule ha fallado y no ha podido cargarse.
         */
        return 0;
}

/**
*	Defino que es lo que se va a hacer al terminar el modulo
*/
static void terminar(void)
{
	remove_proc_entry(FileProc,NULL);
	printk(KERN_INFO "Curso: %s\n", Curso);
}


/*
 * Indicamos cuales son las funciones de inicio y fin
 */
module_init(iniciar);
module_exit(terminar);

/*
 * Documentacion del modulo
 */
MODULE_LICENSE("GPL");
MODULE_AUTHOR("Nery Galvez - 201403525");
MODULE_DESCRIPTION("Modulo con informacion de la memoria RAM");
MODULE_SUPPORTED_DEVICE("TODOS");
