/**
Los estados se extraen de 
/usr/src/linux-headers-4.4.0-127/include/linux/sched.h
#define TASK_RUNNING		0
#define TASK_INTERRUPTIBLE	1
#define TASK_UNINTERRUPTIBLE	2
#define __TASK_STOPPED		4
#define __TASK_TRACED		8
#define EXIT_DEAD		16
#define EXIT_ZOMBIE		32
*/
#include <linux/fs.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/sched.h>
#include <linux/module.h>
#include <linux/seq_file.h>
#include <linux/proc_fs.h>

static int meminfo_proc_show(struct seq_file *m, void *v)
{
    struct task_struct *task;
    int i = 0;
    int zombie = 0, suspendido = 0, detenido = 0, ejecucion = 0;
    seq_printf(m, "{\"processes\":");
    seq_printf(m, "\t[");
    seq_printf(m, "\n");
    for_each_process(task)
    {
        if (i != 0)
            seq_printf(m, ",\n");
        seq_printf(m, "\t\t{\"pid\":%d,\"user\":%d,\"name\":\"%s\",\"ram\":%llu,\"cpu\":%d,\"status\":",
                   task->pid, task->cred->uid, task->comm, task->acct_vm_mem1, task->cpuset_mem_spread_rotor);
        if (task->state == 0)
        {
            seq_printf(m, "\"Ejecutandose\"}");
            ejecucion++;
        }
        if (task->state == 1)
        {
            seq_printf(m, "\"Hibernado\"}");
            suspendido++;
        }
        if (task->state == 2)
        {
            seq_printf(m, "\"U\"}");
            suspendido++;
        }
        if (task->state == 4)
        {
            seq_printf(m, "\"Detenido\"}");
            detenido++;
        }
        if (task->state == 8)
        {
            seq_printf(m, "\"T\"}");
            detenido++;
        }
        if (task->state == 16)
        {
            seq_printf(m, "\"Muerto\"}");
            detenido++;
        }
        if (task->state == 32)
        {
            seq_printf(m, "\"Zombie\"}");
            zombie++;
        }
        i++;
    }
    seq_printf(m, "\n");
    seq_printf(m, "\t],\n");
    seq_printf(m, "\"data\":");
    seq_printf(m, "{");
    seq_printf(m, "\"total\":");
    seq_printf(m, "%d,", i);
    seq_printf(m, "\"ejecucion\":");
    seq_printf(m, "%d,", ejecucion);
    seq_printf(m, "\"suspendidos\":");
    seq_printf(m, "%d,", suspendido);
    seq_printf(m, "\"detenidos\":");
    seq_printf(m, "%d,", detenido);
    seq_printf(m, "\"zombie\":");
    seq_printf(m, "%d", zombie);
    seq_printf(m, "}");
    seq_printf(m, "\n");
    seq_printf(m, "}");
    return 0;
}

static void __exit final(void) //Salida de modulo
{
    printk(KERN_INFO "Cleaning Up.\n");
}

static int meminfo_proc_open(struct inode *inode, struct file *file)
{
    return single_open(file, meminfo_proc_show, NULL);
}

static const struct file_operations meminfo_proc_fops = {
    .open = meminfo_proc_open,
    .read = seq_read,
    .llseek = seq_lseek,
    .release = single_release,
};
static int __init inicio(void) //Escribe archivo en /proc
{
    proc_create("m_processes_201443728", 0, NULL, &meminfo_proc_fops);
    return 0;
}
module_init(inicio);
module_exit(final);
MODULE_LICENSE("GPL");


