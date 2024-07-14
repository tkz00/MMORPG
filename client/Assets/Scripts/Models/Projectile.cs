using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Projectile : MonoBehaviour
{
    [SerializeField]
    private Transform model;

    private Coroutine movementCoroutine;

    private float projectileSpeed = 15f;
    
    public void Move(Vector3 target)
    {
        if (movementCoroutine != null)
        {
            StopCoroutine(movementCoroutine);
        }
        movementCoroutine = StartCoroutine(MoveTowards(new Vector3(target.x, this.transform.position.y, target.z)));
    }
    private IEnumerator MoveTowards(Vector3 target)
    {
        float distanceToGround = transform.position.y - target.y;
        target.y += distanceToGround;
        while (Vector3.Distance(transform.position, target) > 0.1f)
        {
            Vector3 direction = target - transform.position;
            Vector3 movement = direction.normalized * projectileSpeed * Time.deltaTime;
            transform.position += movement;
            transform.rotation = Quaternion.LookRotation(direction.normalized);
            yield return null;
        }
    }

    public void SetScale(float scale)
    {
        model.localScale = Vector3.one * scale;
    }

    public void TriggerHit()
    {
        GetComponent<MeshRenderer>().material.color = Color.red;
    }
}
