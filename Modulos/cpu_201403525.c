/**
*
*	Creando el modulo de cpu
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
#include <linux/sched/signal.h>
#include <linux/sched.h>
 

#define FileProc "cpu_201403525"
#define Carne "201403525"
#define Nombre "Nery Gonzalo Galvez Gomez"
#define SO "Ubuntu 18.04.4 LTS"
#define Curso "Sistemas Operativos 2"

struct task_struct *task;
struct task_struct *task_child;
struct list_head *list;







static int proc_llenar_archivo(struct seq_file *m, void *v) {
        seq_printf(m, "Carne: %s\n", Carne);
        seq_printf(m, "Nombre: %s\n", Nombre);
        seq_printf(m, "Sistema operativo: %s\n", SO);

	//Imprimo la informacion de cada uno de los procesos
	for_each_process(task){
        seq_printf(m, "PID : %d, Nombre : %s, Estado : %ld\n", task->pid, task->comm, task->state);
	list_for_each(list, &task->children){
		task_child = list_entry(list, struct task_struct, sibling);
	        seq_printf(m, "PID : %d, Nombre : %s, Estado : %ld\n", task_child->pid, task_child->comm, task_child->state);
		}
	}
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
	printk(KERN_INFO "Nombre: %s\n", Nombre);
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
MODULE_DESCRIPTION("Modulo con informacion del CPU");
MODULE_SUPPORTED_DEVICE("TODOS");
